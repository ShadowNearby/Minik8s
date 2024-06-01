package activator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/autoscaler"
	"minik8s/pkgs/serverless/function"
	"minik8s/utils"
	"net"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

// GenerateReplicaSet  Generate Replicaset to save function*/
func GenerateReplicaSet(name string, namespace string, image string, replicas int) *core.ReplicaSet {
	return &core.ReplicaSet{
		Kind:       "ReplicaSet",
		ApiVersion: "extensions/v1beta1",
		MetaData: core.MetaData{
			Name:      name,
			Namespace: namespace,
		},
		Spec: core.ReplicaSetSpec{
			Replicas: replicas,
			Selector: core.Selector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: core.ReplicaSetTemplate{
				MetaData: core.MetaData{
					Name:      name,
					Namespace: namespace,
					Labels:    map[string]string{"app": name},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: core.PullAlways,
							Ports: []core.PortConfig{
								{
									ContainerPort: 80,
									Protocol:      "TCP",
									Name:          "p1",
								},
							},
							Cmd: []string{
								"python3",
								"server.py",
							},
						},
					},
				},
			},
		},
		Status: core.ReplicaSetStatus{
			RealReplicas: 0,
			Scale:        0,
			OwnerReference: core.OwnerReference{
				Name:       name,
				ObjType:    core.ObjFunction,
				Controller: true,
			},
		},
	}
}

func getPodIpList(pods []core.Pod) []string {
	result := make([]string, 0)
	if pods == nil {
		return result
	}
	for _, pod := range pods {
		if pod.Status.Phase != core.PhasePending && pod.Status.PodIP != "" {
			result = append(result, pod.Status.PodIP)
		}
	}
	log.Info(result)
	return result
}
func CheckConnection(ip string) error {
	timer := time.NewTimer(config.FunctionConnectTime)
	for {
		select {
		case <-timer.C:
			return errors.New("timeout")
		default:
			{
				// try to connect to the ip
				address := ip + ":" + config.ServerlessIP
				conn, err := net.DialTimeout("tcp", address, time.Second)
				if err != nil {
					continue
				}
				defer conn.Close()
				log.Info("[CheckConnection] Connection is ok")
				return nil
			}
		}

	}
}

func InitFunc(name string, path string) error {
	err := function.CreateImage(path, name)
	if err != nil {
		log.Error("[InitFunc] create image error: ", err)
		return err
	}
	imageName := fmt.Sprintf("%s:%s/%s:latest", config.LocalServerIp, config.ApiServerPort, name)

	replicaSet := GenerateReplicaSet(name, "serverless", imageName, 0)
	log.Info("[InitFunc] create record replicaSet: ", replicaSet)

	// create the record
	log.Info("[InitFunc] create the record")
	autoscaler.RecordMutex.Lock()
	autoscaler.RecordMap[name] = &autoscaler.Record{
		Name:      name,
		Replicas:  0,
		PodIps:    make(map[string]int32),
		CallCount: 0,
	}
	autoscaler.RecordMutex.Unlock()
	log.Info("[InitFunc] create the record successfully")
	err = utils.CreateObject(core.ObjReplicaSet, replicaSet.MetaData.Namespace, replicaSet)
	if err != nil {
		log.Error("[InitFunc] create record error: ", err)
	}
	return nil
}
func IfDeployed(name string) ([]string, error) {
	log.Info("[CheckPrepare] check prepare for function: ", name)
	// 1. find the pods

	response := utils.GetObject(core.ObjReplicaSet, "serverless", name)
	var replicaSet *core.ReplicaSet
	if response == "" {
		return nil, errors.New("cannot get replica")
	}
	err := utils.JsonUnMarshal(response, &replicaSet)
	if err != nil {
		log.Error("[CheckPrepare] error unmarshalling replicas: ", err)
		return nil, err
	}
	log.Info("[CheckPrepare] check prepare for function: ", name)
	retry := config.FunctionRetryTimes
	firstTry := true
	for retry > 0 {
		timer := time.NewTimer(2 * config.FunctionConnectTime)
		deployed := false
		retry--
		for {
			select {
			case <-timer.C:
				log.Info("[CheckPrepare] check prepare for function: ", name)
				break
			default:
				if !deployed {
					log.Info("[CheckPrepare] first check if function deployed: ", name)
					pods, _ := utils.FindRSPods(replicaSet.MetaData.Name)
					podIps := getPodIpList(pods)
					autoscaler.RecordMutex.Lock()
					record := autoscaler.GetRecord(name)
					if record == nil {
						autoscaler.RecordMap[name] = &autoscaler.Record{
							Name:      name,
							Replicas:  int32(len(podIps)),
							PodIps:    make(map[string]int32),
							CallCount: 1,
						}
						record = autoscaler.RecordMap[name]
					} else {
						if firstTry {
							record.CallCount++
							firstTry = false
						}
						record.Replicas = int32(len(podIps))
						autoscaler.RecordMap[name] = record
					}
					log.Info(pods, podIps)
					if record.CallCount > replicaSet.Status.Scale && record.CallCount < config.FunctionThreshold {
						replicaSet.Status.Scale = record.CallCount
						log.Info("[CheckPrepare] scale up the function: ", name, "the replica number: ", replicaSet.Status.Scale)
						err = utils.UpdateObject(core.ObjReplicaSet, replicaSet.MetaData.Namespace, name, replicaSet)
						if err != nil {
							log.Error("[CheckPrepare] update replica set error: ", err)
						}
					} else {
						autoscaler.RecordMutex.Unlock()
						if len(podIps) > 0 {
							return podIps, nil
						}
					}
					autoscaler.RecordMutex.Unlock()
					deployed = true
				} else {
					pods, _ := utils.FindRSPods(replicaSet.MetaData.Name)
					autoscaler.RecordMutex.RLock()
					record := autoscaler.GetRecord(name)
					autoscaler.RecordMutex.RUnlock()
					if record == nil {
						log.Error("[CheckPrepare] record not found")
						return nil, errors.New("record not found")
					}
					podsIp := getPodIpList(pods)
					log.Info("[CheckPrepare] the pod ip list in second time or later: ", podsIp)
					log.Info("the replica number: ", int32(len(podsIp)), " the call count: ", record.CallCount)
					if int32(len(podsIp)) >= record.CallCount && len(podsIp) > 0 {
						// update the replica count
						record.Replicas = int32(len(pods))
						autoscaler.RecordMutex.Lock()
						autoscaler.RecordMap[name] = record
						autoscaler.RecordMutex.Unlock()
						log.Info("[CheckPrepare] the replica number is correct: ", int32(len(podsIp)), record.CallCount)
						return podsIp, nil
					}
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
	// get the current pod ip list and return
	log.Info("[CheckPrepare] get the current pod ip list and return")
	pods, _ := utils.FindRSPods(replicaSet.MetaData.Name)
	podsIp := getPodIpList(pods)
	return podsIp, nil
}

// DeleteFunc delete the function
func DeleteFunc(name string) error {

	// 1. delete the replicaset
	err := utils.DeleteObject(core.ObjReplicaSet, "serverless", name)

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
	// 3. delete the old image
	err = function.DeleteImage(name)
	if err != nil {
		log.Error("[DeleteFunc] delete image error: ", err)
		return err
	}
	return nil
}

// TriggerFunc trigger the function with some parameters(retry 3 times)
// if the function is not deployed, deploy it first
func TriggerFunc(name string, params []byte) ([]byte, error) {
	// 1. check if the function is deployed
	log.Info("[TriggerFunc] trigger function: ", name)
	retry := 3
	for retry > 0 {
		retry -= 1
		podIps, err := IfDeployed(name)
		if err != nil {
			log.Error("[TriggerFunc] check prepare error: ", err)
			continue
		}
		// 2. load balance
		podIp, err := LoadBalance(name, podIps)
		if err != nil {
			log.Error("[TriggerFunc] load balance error: ", err)
			continue
		}
		// 3. trigger the function
		log.Info("[TriggerFunc] load balance success: ", podIp)
		url := fmt.Sprintf("http://%s:18080/", podIp)
		if err != nil {
			return nil, err
		}
		var data interface{}
		err = json.Unmarshal(params, &data)
		prettyJSON, err := json.MarshalIndent(data, "", "  ")
		err = utils.JsonUnMarshal(string(params), &data)
		if err != nil {
			log.Error("[TriggerFunc] marshal params error: ")
			{
				return nil, err
			}
		}

		log.Info("[TriggerFunc] prettyJSON: ", string(prettyJSON), "url: ", url)

		// 4. send the request
		// first check the connection
		err = CheckConnection(podIp)
		if err != nil {
			log.Error("[TriggerFunc] check connection error: ", err)
			continue
		}
		log.Info("[TriggerFunc] connection is finished")
		_, ret, err := utils.SendRequest("POST", url, params)
		log.Info("[TriggerFunc] ret: ", string(ret))
		result := bytes.NewBufferString(ret).Bytes()
		if err != nil {
			log.Error("[TriggerFunc] send request error: ", err)
			continue
		}
		return result, err
	}
	return nil, errors.New("trigger function error")
}

// LoadBalance choose a pod ip to trigger the function
func LoadBalance(name string, podIps []string) (string, error) {
	if len(podIps) == 0 {
		log.Error("[LoadBalance] pod ip list is empty")
		return "", errors.New("pod ip list is empty")
	}

	autoscaler.RecordMutex.RLock()
	record := autoscaler.GetRecord(name)
	autoscaler.RecordMutex.RUnlock()

	if record == nil {
		log.Error("[LoadBalance] record not found")
		return "", errors.New("record not found")
	}

	// update the record
	for _, podIp := range podIps {
		if _, ok := record.PodIps[podIp]; !ok {
			record.PodIps[podIp] = 0
		}
	}
	// choose the pod ip with the least call count
	sort.Slice(podIps, func(i, j int) bool {
		return record.PodIps[podIps[i]] < record.PodIps[podIps[j]]
	})
	chosenPodIp := podIps[0]
	record.PodIps[chosenPodIp]++

	autoscaler.RecordMutex.Lock()
	autoscaler.RecordMap[name] = record
	autoscaler.RecordMutex.Unlock()

	return chosenPodIp, nil
}
