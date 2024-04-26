package test

import (
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/resources"
	"minik8s/utils"
)

func MetricsTest() {
	podConfig := GeneratePodConfigPy()
	err := controller.CreatePod(&podConfig)
	//err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		logger.Errorf("cannot create pod")
		resources.StopPod(podConfig)
		return
	}
	logger.Infof("-------------CREATE POD FINISHED-----------")
	stat := controller.NodeMetrics()
	text := utils.JsonMarshal(stat)
	logger.Infof("%s", text)
	//ms, err := resources.GetPodMetrics(&podConfig)
	//if err != nil {
	//	logger.Errorf("cannot inspect pod")
	//}
	//for i := range ms {
	//	text := utils.JsonMarshal(ms[i])
	//	logger.Infof("%s", text)
	//}
	resources.StopPod(podConfig)
}
