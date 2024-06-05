package autoscaler

import (
	"errors"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Record struct {
	/* Name is the name of the function */
	Name string `json:"name"`
	/* current replica number of the function */
	Replicas int `json:"replicas"`
	// the podIps that the function has deployed on
	PodIps map[string]int `json:"podIps"`
	// the call count of the function
	CallCount int `json:"callCount"`
	// last call time
	LastCallTime time.Time `json:"lastCallTime"`
}

var (
	RecordMap   = make(map[string]Record)
	RecordMutex sync.RWMutex // protect the access of RecordMap
)

func RecordBackGroundCheck() {
	ticker := time.NewTicker(config.ServerlessScaleToZeroTime)
	defer ticker.Stop()
	for range ticker.C {
		RecordMutex.Lock()
		currTime := time.Now()
		for key, record := range RecordMap {
			if currTime.Sub(record.LastCallTime) > config.ServerlessScaleToZeroTime {
				log.Info("rescale to zero")
				// should rescale to zero
				// get the replicaset and reset the spec.replica to zero
				var replica core.ReplicaSet
				rsTxt := utils.GetObject(core.ObjReplicaSet, "default", record.Name)
				utils.JsonUnMarshal(rsTxt, &replica)
				replica.Spec.Replicas = 0
				utils.SetObject(core.ObjReplicaSet, "default", record.Name, replica)
				record.CallCount = 0
				RecordMap[key] = record
			}
		}
		RecordMutex.Unlock()
	}
}

func GetRecord(name string) (Record, error) {
	if record, ok := RecordMap[name]; ok {
		return record, nil
	}
	return Record{}, errors.New("no record available")
}

func SetRecord(name string, record Record) {
	record.LastCallTime = time.Now()
	RecordMap[name] = record
}

func DeleteRecord(name string) {
	delete(RecordMap, name)
}

func UpdateRecord(name string) {
	if record, ok := RecordMap[name]; ok {
		record.CallCount++
		record.LastCallTime = time.Now()
		RecordMap[name] = record
	}
}

// LoadBalance choose a pod ip to trigger the function, will maintain metadata at the same time
func LoadBalance(name string, podIps []string) (string, error) {
	if len(podIps) == 0 {
		log.Error("[LoadBalance] pod ip list is empty")
		return "", errors.New("pod ip list is empty")
	}

	RecordMutex.RLock()
	record, err := GetRecord(name)
	RecordMutex.RUnlock()
	if err != nil {
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

	RecordMutex.Lock()
	RecordMap[name] = record
	RecordMutex.Unlock()

	return chosenPodIp, nil
}
