package autoscaler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"time"
)

// query the pod ips of the replicaSet
func QueryPodIps() (map[string][]string, error) {
	// find all the pods of the replicaSet
	response := utils.GetObject(core.ObjPod, "serverless", "")
	result := make(map[string][]string)
	var podList []core.Pod
	err := json.Unmarshal([]byte(response), &podList)
	if err != nil {
		log.Error("[QueryPodIps] error unmarshalling pods: ", err)
		return nil, err
	}
	for _, pod := range podList {
		pos, ok := result[pod.MetaData.Name]
		if !ok {
			pos = make([]string, 0)
		}
		pos = append(pos, pod.Status.PodIP)
		result[pod.MetaData.Name] = pos
	}
	return result, nil
}

// PeriodicMetric check the invoke frequency periodically,
// // delete the function if it is not invoked for a long time
func PeriodicMetric(timeInterval int) {
	for {
		response := utils.GetObject(core.ObjReplicaSet, "serverless", "")
		// get all replicas

		replicaList := &[]core.ReplicaSet{}
		err := json.Unmarshal([]byte(response), replicaList)
		if err != nil {
			log.Error("[PeriodicMetric] error unmarshalling replicas: ", err)
			continue
		}
		// update the replicas information
		for _, replica := range *replicaList {
			// get the according record in map
			RecordMutex.RLock()
			record := GetRecord(replica.MetaData.Name)
			RecordMutex.RUnlock()
			if record == nil {
				//TODO:Ready ? Real?
				record = &Record{
					Name:      replica.MetaData.Name,
					Replicas:  replica.Status.ReadyReplicas,
					PodIps:    make(map[string]int32),
					CallCount: 0,
				}
				RecordMutex.Lock()
				SetRecord(replica.MetaData.Name, record)
				RecordMutex.Unlock()
			} else {
				// if the call times is 0, scale to zero
				// scale according to the call times
				replica.Status.Scale = record.CallCount
				// update the replicaset
				if replica.Status.Scale != replica.Status.ReadyReplicas {
					err := utils.UpdateObject(core.ObjReplicaSet, replica.MetaData.Name, replica)
					if err != nil {
						log.Error("[PeriodicMetric] error updating replicas: ", err)
						continue
					}
				}

				// the replica count in expection
				record.Replicas = record.CallCount
				record.CallCount = 0
				RecordMutex.Lock()
				SetRecord(replica.MetaData.Name, record)
				RecordMutex.Unlock()
			}
		}
		time.Sleep(time.Duration(timeInterval) * time.Second)
	}
}
