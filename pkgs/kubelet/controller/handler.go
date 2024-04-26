package controller

import (
	"context"
	"errors"
	"github.com/containerd/containerd"
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
	pStatChan := make(chan core.PodStatus, 2)
	ctNameChan := make(chan resources.NameIdPair, cLen)
	doneChan := make(chan bool)
	runtime.KubeletInstance.WritePodConfig(pConfig.MetaData.Name, pConfig.MetaData.NameSpace, pConfig)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup, dChan chan<- bool) {
		defer wg.Done()
		err := resources.CreatePod(pConfig, pStatChan, ctNameChan, dChan)
		if err != nil {
			doneChan <- false
			logger.Errorf("create pod error: %s", err.Error())
		}
	}(&wg, doneChan)
	// pod status
	var pStat core.PodStatus
	for i := 0; i < cLen+2; i++ {
		select {
		case pStat = <-pStatChan:
			logger.Infof("init pod status")
		case ctNameID := <-ctNameChan:
			runtime.KubeletInstance.ContainerStart(&pStat, ctNameID.Name, ctNameID.ID)
		case done := <-doneChan:
			if done == true {
				// pod init done
				pStat.Status = core.PhaseRunning
				break
			} else {
				return errors.New("create pod failed")
			}
		case <-time.After(10 * time.Second):
			return errors.New("create pod time out: 10 secs")
		}
	}
	close(pStatChan)
	close(ctNameChan)
	close(doneChan)
	wg.Wait()
	runtime.KubeletInstance.WritePodStat(pConfig.MetaData.Name, pConfig.MetaData.NameSpace, &pStat)
	return nil
}

// StopPod stop and remove container
func StopPod(pConfig core.Pod) error {
	runtime.KubeletInstance.DelPodConfig(pConfig.MetaData.Name, pConfig.MetaData.NameSpace)
	runtime.KubeletInstance.DelPodStat(pConfig.MetaData.Name, pConfig.MetaData.NameSpace)
	return resources.StopPod(pConfig)
}

// InspectPod exec_probe of the pod, if a pod failed, then stop it
func InspectPod(pConfig core.Pod, probeType runtime.ProbeType) string {
	containers := resources.ContainerManagerInstance.GetPodContainers(&pConfig)
	containerMap := make(map[string]containerd.Container, len(containers))
	for _, container := range containers {
		id := container.ID()
		image, err := container.Image(context.Background())
		if err != nil {
			continue
		}
		if image.Name() == "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9" {
			continue
		}
		containerMap[id] = container
	}
	err := runtime.KubeletInstance.DoProbe(probeType, containerMap, pConfig)
	//err := runtime.KubeletInstance.LivenessProbe(containerMap, pConfig)
	if err != nil {
		logger.Errorf("liveness probe error: %s", err.Error())
		return ""
	}
	// print status
	pStat := runtime.KubeletInstance.GetPodStat(pConfig.MetaData.Name, pConfig.MetaData.NameSpace)
	jsonText := utils.CreateJson(pStat)
	logger.Infof("live probe:\n%s", jsonText)
	return jsonText
}

// NodeMetrics return the metrics of a node, including ready, cpu, memory, process_num, disk, network
func NodeMetrics() core.NodeMetrics {
	var allPID uint64
	var allMem uint64
	var allCPU uint64
	var allDisk uint64
	logger.Infof("len:%d", len(runtime.KubeletInstance.PodConfigMap))
	for name, podConfig := range runtime.KubeletInstance.PodConfigMap {
		metrics, err := resources.GetPodMetrics(&podConfig)
		if err != nil {
			logger.Errorf("get pod %s metrics error: %s", name, err.Error())
			continue
		}
		for _, metric := range metrics {
			allPID += metric.PidCount
			allMem += metric.MemoryUsage
			allCPU += metric.CpuUsage
			allDisk += metric.DiskUsage
		}
	}
	return core.NodeMetrics{
		Ready:              true,
		CPUUsage:           allCPU,
		MemoryUsage:        allMem,
		PIDUsage:           allPID,
		DiskUsage:          allDisk,
		NetworkUnavailable: false,
	}
}
