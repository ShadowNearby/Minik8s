package test

import (
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/kubelet/resources"
)

func MetricsTest() {
	podConfig := GeneratePodConfigPy()
	err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		logger.Errorf("cannot create pod")
		resources.StopPod(podConfig)
		return
	}
	logger.Infof("-------------CREATE POD FINISHED-----------")
	err = resources.InspectPod(&podConfig)
	if err != nil {
		logger.Errorf("cannot inspect pod")
	}
	resources.StopPod(podConfig)
}
