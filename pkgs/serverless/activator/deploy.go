package activator

import (
	"errors"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/controller/autoscaler"
	"minik8s/pkgs/serverless/function"
	"minik8s/utils"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// InitFunction will create an image and registry it in master node
// then generate a replicaset to manage pods
func InitFunction(name string, path string) error {
	err := function.CreateImage(path, name)
	if err != nil {
		log.Error("[InitFunction] create image error: ", err)
		return err
	}
	imageName := fmt.Sprintf("%s:5000/%s:v1", config.ClusterMasterIP, name)
	replicaSet := utils.GenerateRSConfig(name, "default", imageName, 0)
	log.Info("[InitFunction] create record replicaSet: ", replicaSet)

	// create the record
	log.Info("[InitFunction] create the record")
	autoscaler.RecordMutex.Lock()
	autoscaler.RecordMap[name] = autoscaler.Record{
		Name:      name,
		Replicas:  0,
		PodIps:    make(map[string]int),
		CallCount: 0,
	}
	autoscaler.RecordMutex.Unlock()
	log.Infof("create replicaset %s", replicaSet.MetaData.Name)
	err = utils.CreateObject(core.ObjReplicaSet, "default", replicaSet)
	if err != nil {
		log.Errorf("create rs failed: %s", err.Error())
	}
	return nil
}

// DeleteFunc delete the function, including replicaset, record and docker image
func DeleteFunc(name string) error {
	// 1. delete the replicaset
	err := utils.DeleteObject(core.ObjReplicaSet, "default", name)

	if err != nil {
		log.Error("[DeleteFunc] delete replicas error: ", err)
		return err
	}

	// 2. delete the record from the record map
	log.Info("[DeleteFunc] delete record from record map")
	autoscaler.RecordMutex.Lock()
	autoscaler.DeleteRecord(name)
	autoscaler.RecordMutex.Unlock()

	log.Info("[DeleteFunc] delete record from record map")

	// 3. delete the image
	err = function.DeleteImage(name)
	if err != nil {
		log.Error("[DeleteFunc] delete image error: ", err)
		return err
	}
	return nil
}

// TriggerFunc trigger the function with some parameters
// if the function is not deployed, deploy it first
func TriggerFunc(name string, params []byte) (string, error) {
	// 1. check if the function is deployed
	podIps, err := getAvailablePods(name)
	if err != nil {
		log.Error("[TriggerFunc] check prepare error: ", err)
		return "", errors.New("cannot asign pod to node")
	}
	if len(podIps) == 0 {
		return "", errors.New("cannot get available pod")
	}
	// 2. load balance
	podIp, err := autoscaler.LoadBalance(name, podIps)
	if err != nil {
		log.Error("[TriggerFunc] load balance error: ", err)
		return "", errors.New("cannot load balance")
	}

	// 3. trigger the function
	url := fmt.Sprintf("http://%s:18080", podIp)
	err = checkConnection(podIp)
	if err != nil {
		log.Error("[TriggerFunc] check connection error: ", err)
		return "", errors.New("cannot connect to selected node")
	}
	request := core.TriggerRequest{
		Url:    url,
		Params: params,
	}
	info, err := utils.SendTriggerRequest(request)
	if err != nil {
		log.Errorf("[SendTriggerRequest] tigger request failed: %s", err.Error())
		return "", err
	}
	return info, nil
}

func getAvailablePods(name string) ([]string, error) {

	for i := 0; i < config.FunctionRetryTimes; i++ {
		autoscaler.RecordMutex.Lock()
		replicaSet, err := utils.FindFunctionRs(name)
		if err != nil {
			log.Errorf("cannot find serverless replicaset: %s", err.Error())
			return nil, err
		}
		pods, err := utils.FindRSPods(true, replicaSet.MetaData.Name, replicaSet.MetaData.Namespace)
		if err != nil {
			log.Errorf("cannot find rs's pods: %s", err.Error())
		}
		log.Infof("find rs pods: %d", len(pods))
		podIps := getPodIpList(&pods)
		record, err := autoscaler.GetRecord(name)
		if err != nil {
			autoscaler.RecordMap[name] = autoscaler.Record{
				Name:         name,
				Replicas:     len(podIps),
				PodIps:       make(map[string]int),
				CallCount:    1,
				LastCallTime: time.Now(),
			}
			record = autoscaler.RecordMap[name]
		} else {
			log.Infof("call count: %d", record.CallCount)
			record.CallCount++
			// record.Replicas = len(podIps)
			record.LastCallTime = time.Now()
			autoscaler.RecordMap[name] = record
		}
		expectReplica := (record.CallCount + 9) / 10
		log.Infof("expect replica: %d, replica: %d", expectReplica, record.Replicas)
		if expectReplica > record.Replicas && record.Replicas < config.FunctionThreshold {
			replicaSet.Spec.Replicas = expectReplica
			record.Replicas = expectReplica
			autoscaler.RecordMap[name] = record
			log.Infof("call count: %d, real replica: %d, scale up %s to %d", record.CallCount, replicaSet.Status.RealReplicas, name, replicaSet.Spec.Replicas)
			err = utils.SetObject(core.ObjReplicaSet, replicaSet.MetaData.Namespace, replicaSet.MetaData.Name, replicaSet)
			if err != nil {
				log.Error("[CheckPrepare] update replica set error: ", err)
			}
		} else {
			if len(podIps) > 0 {
				autoscaler.RecordMutex.Unlock()
				return podIps, nil
			}
		}
		autoscaler.RecordMutex.Unlock()
		if len(podIps) > 0 {
			return podIps, nil
		}
		time.Sleep(3 * time.Second)
	}
	return nil, errors.New("cannot get availabel pods")

}

func getPodIpList(pods *[]core.Pod) []string {
	result := make([]string, 0)
	if pods == nil {
		return result
	}
	for _, pod := range *pods {
		log.Infof("phase: %s, podIP: %s", pod.Status.Phase, pod.Status.PodIP)
		if pod.Status.Phase == core.PodPhaseRunning && pod.Status.PodIP != "" {
			result = append(result, pod.Status.PodIP)
		}
	}
	return result
}
func checkConnection(ip string) error {
	timer := time.NewTimer(config.FunctionConnectTime)
	for {
		select {
		case <-timer.C:
			return errors.New("timeout")
		default:
			{
				address := fmt.Sprintf("%s:%s", ip, "18080") // TODO: cannot read config
				conn, err := net.DialTimeout("tcp", address, time.Second)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}
				defer conn.Close()
				return nil
			}
		}

	}
}
