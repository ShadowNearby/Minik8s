package kubeletcontroller

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/resources"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"sync"
	"time"
)

// CreatePod pull and create containers of a pod, and register the pod to kubelet runtime
func CreatePod(pConfig *core.Pod) error {
	cLen := len(pConfig.Spec.Containers)
	pStatusChan := make(chan core.PodStatus)
	cStatusChan := make(chan core.ContainerStatus, cLen)
	doneChan := make(chan bool)
	runtime.KubeletInstance.WritePodConfig(pConfig.MetaData.Name, pConfig.MetaData.Namespace, pConfig)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup, dChan chan<- bool) {
		defer wg.Done()
		err := resources.CreatePod(pConfig, pStatusChan, cStatusChan, dChan)
		if err != nil {
			doneChan <- false
			logger.Errorf("create pod error: %s", err.Error())
		}
	}(&wg, doneChan)
	// pod status
	var pStat core.PodStatus
	for i := 0; i < cLen+3; i++ {
		select {
		case pStat = <-pStatusChan:
			logger.Infof("init pod status")
		case cStatus := <-cStatusChan:
			pStat.ContainersStatus = append(pStat.ContainersStatus, cStatus)
		case done := <-doneChan:
			if done == true {
				pStat.StartTime = time.Now()
				pStat.Condition = core.ConReady
				break
			} else {
				return errors.New("create pod failed")
			}
		case <-time.After(30 * time.Second):
			close(pStatusChan)
			close(cStatusChan)
			close(doneChan)
			return errors.New("create pod time out: 10 secs")
		}
	}
	wg.Wait()
	pConfig.Status = pStat
	runtime.KubeletInstance.WritePodConfig(pConfig.MetaData.Name, pConfig.MetaData.Namespace, pConfig)
	return nil
}

// StopPod stop and remove container
func StopPod(pConfig core.Pod) error {
	runtime.KubeletInstance.DelPodConfig(pConfig.MetaData.Name, pConfig.MetaData.Namespace)
	runtime.KubeletInstance.DelPodStat(pConfig.MetaData.Name, pConfig.MetaData.Namespace)
	return resources.StopPod(pConfig)
}

// InspectPod exec_probe of the pod, if a pod failed, then stop it
func InspectPod(pod *core.Pod, probeType runtime.ProbeType) string {
	containers := resources.ContainerManagerInstance.GetPodContainers(pod)
	logger.Infof("container len: %d", len(containers))
	err := runtime.KubeletInstance.DoProbe(probeType, containers, pod)
	if err != nil {
		logger.Errorf("liveness probe error: %s", err.Error())
		return ""
	}
	// print status
	pStat := pod.Status
	jsonText := utils.JsonMarshal(pStat)
	logger.Infof("live probe:\n%s", jsonText)
	return jsonText
}

// NodeMetrics return the metrics of a node, including ready, cpu, memory, process_num, disk, network
func NodeMetrics() core.NodeMetrics {
	state, err := runtime.GetNodeState()
	if err != nil {
		logger.Errorf("error: %s", err.Error())
		return state
	}
	return state
}

func PodMetrics(pod core.Pod) core.Metrics {
	metrics, err := runtime.GetPodMetrics(&pod)
	if err != nil {
		logger.Errorf("get pod metrics error")
		return core.Metrics{}
	}
	return metrics
}
