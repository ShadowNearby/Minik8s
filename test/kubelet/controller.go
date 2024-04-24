package test

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"time"
)

func CreatePodTest() {
	podConfig := GeneratePodConfigPy()
	err := controller.CreatePod(&podConfig)
	if err != nil {
		logger.Errorf("controller create pod error: %s", err.Error())
		return
	}
	pStatus := runtime.KubeletInstance.PodMap[podConfig.MetaData.Name]
	marshaledJSON, err := json.MarshalIndent(pStatus, "", "  ")
	if err != nil {
		logger.Errorf("error: %s", err.Error())
		_ = controller.StopPod(podConfig)
		return
	}
	logger.Infof("inspect data: %s", marshaledJSON)
	err = controller.StopPod(podConfig)
	if err != nil {
		logger.Errorf("stop pod error: %s", err.Error())
	}
}

func InspectPod() {
	podConfig := GeneratePodConfigPy()
	err := controller.CreatePod(&podConfig)
	if err != nil {
		return
	}
	controller.InspectPod(podConfig, runtime.ExecProbe)

	time.Sleep(6 * time.Second)

	controller.InspectPod(podConfig, runtime.ExecProbe)

	controller.StopPod(podConfig)
}
