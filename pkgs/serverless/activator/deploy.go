package activator

import (
	"errors"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/autoscaler"
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
	imageName := fmt.Sprintf("shadownearby/%s:v1", name)
	replicaSet := utils.GenerateRSConfig(name, "default", imageName, 0)
	log.Info("[InitFunction] create record replicaSet: ", replicaSet)

	// create the record
	log.Info("[InitFunction] create the record")
	autoscaler.RecordMutex.Lock()
	autoscaler.RecordMap[name] = &autoscaler.Record{
		Name:      name,
		Replicas:  0,
		PodIps:    make(map[string]int),
		CallCount: 0,
	}
	autoscaler.RecordMutex.Unlock()
	log.Info("[InitFunction] create the record successfully")
	err = utils.CreateObject(core.ObjReplicaSet, replicaSet.MetaData.Namespace, replicaSet)
	if err != nil {
		log.Error("[InitFunction] create record error: ", err)
	}
	return nil
}

// DeleteFunc delete the function
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
func TriggerFunc(name string, params []byte) error {
	// 1. check if the function is deployed
	podIps, err := getAvailablePods(name)
	if err != nil {
		log.Error("[TriggerFunc] check prepare error: ", err)
		return errors.New("cannot asign pod to node")
	}
	// 2. load balance
	podIp, err := autoscaler.LoadBalance(name, podIps)
	if err != nil {
		log.Error("[TriggerFunc] load balance error: ", err)
		return errors.New("cannot load balance")
	}

	// 3. trigger the function
	url := fmt.Sprintf("http://%s:18080", podIp)
	err = checkConnection(podIp)
	if err != nil {
		log.Error("[TriggerFunc] check connection error: ", err)
		return errors.New("cannot connect to selected node")
	}
	request := core.TriggerRequest{
		Url:    url,
		Params: params,
	}
	err = utils.SendTriggerRequest(request)
	if err != nil {
		log.Errorf("[SendTriggerRequest] tigger request failed: %s", err.Error())
		return err
	}
	return nil
}

func getAvailablePods(name string) ([]string, error) {
	replicaSet, err := utils.FindFunctionRs(name)
	if err != nil {
		log.Errorf("cannot find serverless replicaset: %s", err.Error())
		return nil, err
	}
	pods, err := utils.FindRSPods(replicaSet.MetaData.Name, replicaSet.MetaData.Namespace)
	if err != nil {
		log.Errorf("cannot find rs's pods: %s", err.Error())
	}
	podIps := getPodIpList(&pods)
	autoscaler.RecordMutex.Lock()
	record := autoscaler.GetRecord(name)
	if record == nil {
		autoscaler.RecordMap[name] = &autoscaler.Record{
			Name:      name,
			Replicas:  len(podIps),
			PodIps:    make(map[string]int),
			CallCount: 1,
		}
		record = autoscaler.RecordMap[name]
	} else {
		record.CallCount++
		record.Replicas = len(podIps)
		autoscaler.RecordMap[name] = record
	}
	if record.CallCount > replicaSet.Status.RealReplicas && record.CallCount < config.FunctionThreshold {
		replicaSet.Spec.Replicas = record.CallCount
		log.Infof("scale up %s to %d", name, replicaSet.Spec.Replicas)
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

	time.Sleep(10 * time.Second)
	// get the current pod ip list and return
	var podsIp []string
	for i := 0; i < config.FunctionRetryTimes; i++ {
		log.Info("[CheckPrepare] get the current pod ip list and return")
		pods, err = utils.FindRSPods(replicaSet.MetaData.Name, "default")
		if err != nil {
			log.Errorf("find rs pods failed %s", err.Error())
			return nil, err
		}
		podsIp = getPodIpList(&pods)
		if len(podsIp) >= record.CallCount {
			break
		}
		time.Sleep(5 * time.Second)
	}

	return podsIp, nil
}

func getPodIpList(pods *[]core.Pod) []string {
	result := make([]string, 0)
	if pods == nil {
		return result
	}
	for _, pod := range *pods {
		log.Infof("phase: %s, podIP: %s", pod.Status.Condition, pod.Status.PodIP)
		if pod.Status.Condition == core.CondRunning && pod.Status.PodIP != "" {
			log.Info("append into result")
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
				log.Infof("address: %s", address)
				conn, err := net.DialTimeout("tcp", address, time.Second)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}
				defer conn.Close()
				log.Info("[checkConnection] Connection is ok")
				return nil
			}
		}

	}
}
