package test

import (
	"minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"testing"
	"time"
)

func TestInspectPod(t *testing.T) {
	podConfig := GeneratePodConfigPy()
	err := controller.CreatePod(&podConfig)
	if err != nil {
		_ = controller.StopPod(podConfig)
		t.Errorf("cannot create pod")
		return
	}
	controller.InspectPod(&podConfig, runtime.ExecProbe)

	time.Sleep(6 * time.Second)

	controller.InspectPod(&podConfig, runtime.ExecProbe)

	_ = controller.StopPod(podConfig)
}

func TestMetricsTest(t *testing.T) {
	podConfig := GeneratePodConfigPy()
	err := controller.CreatePod(&podConfig)
	//err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot create pod")
		controller.StopPod(podConfig)
		return
	}
	_ = controller.NodeMetrics()
	//text := utils.JsonMarshal(stat)
	controller.StopPod(podConfig)
}
