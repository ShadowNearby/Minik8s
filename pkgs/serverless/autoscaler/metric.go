package autoscaler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"time"
)

// PeriodicMetric check the invoke frequency periodically,
// // delete the function if it is not invoked for a long time
func PeriodicMetric(timeInterval int) {
	for {
		response := utils.GetObject(core.ObjReplicaSet, "serverless", "")
		// get all replicas
		replicaList := &[]core.ReplicaSet{}
		err := json.Unmarshal([]byte(response), replicaList)
		log.Info(response)

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
