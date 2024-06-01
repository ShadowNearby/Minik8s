package autoscaler

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sort"
	"sync"
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
}

var (
	RecordMap   = make(map[string]*Record)
	RecordMutex sync.RWMutex // protect the access of RecordMap
)

func GetRecord(name string) *Record {
	return RecordMap[name]
}

func SetRecord(name string, record *Record) {
	RecordMap[name] = record
}

func DeleteRecord(name string) {
	delete(RecordMap, name)
}

func UpdateRecord(name string) {
	record := GetRecord(name)
	if record == nil {
		return
	}
	record.CallCount++
	SetRecord(name, record)
}

// LoadBalance choose a pod ip to trigger the function, will maintain metadata at the same time
func LoadBalance(name string, podIps []string) (string, error) {
	if len(podIps) == 0 {
		log.Error("[LoadBalance] pod ip list is empty")
		return "", errors.New("pod ip list is empty")
	}

	RecordMutex.RLock()
	record := GetRecord(name)
	RecordMutex.RUnlock()

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

	RecordMutex.Lock()
	RecordMap[name] = record
	RecordMutex.Unlock()

	return chosenPodIp, nil
}
