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

func CreatePod(pConfig *core.Pod) error {
	cLen := len(pConfig.Spec.Containers)
	pStatChan := make(chan core.PodStatus, 2)
	ctNameChan := make(chan resources.NameIdPair, cLen)
	doneChan := make(chan bool)
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
	runtime.KubeletInstance.InsertMap(pConfig.MetaData.Name, pStat)
	return nil
}

func StopPod(pConfig core.Pod) error {
	return resources.StopPod(pConfig)
}

func InspectPod(pConfig core.Pod, probeType runtime.ProbeType) {
	containers := resources.ContainerManagerInstance.GetPodContainers(&pConfig)
	logger.Infof("container len: %d", len(containers))
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
		return
	}
	// print status
	pStat := runtime.KubeletInstance.PodMap[pConfig.MetaData.Name]
	jsonText := utils.JSONPrint(pStat)
	logger.Infof("live probe:\n%s", jsonText)
}