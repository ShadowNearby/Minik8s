package test

import (
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/kubelet"
)

func MetricsTest() {
	podConfig := GeneratePodConfigPy()
	err := kubelet.CreatePod(&podConfig)
	if err != nil {
		logger.Errorf("cannot create pod")
		kubelet.StopPod(&podConfig)
		return
	}
	logger.Infof("-------------CREATE POD FINISHED-----------")
	err = kubelet.InspectPod(&podConfig)
	if err != nil {
		logger.Errorf("cannot inspect pod")
	}
	kubelet.StopPod(&podConfig)
}
