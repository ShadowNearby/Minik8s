package test

import (
	"minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"testing"
	"time"
)

func TestInspectPod(t *testing.T) {
	podConfig := GeneratePodConfigPy()
	err := kubeletcontroller.CreatePod(&podConfig)
	if err != nil {
		_ = kubeletcontroller.StopPod(podConfig)
		t.Errorf("cannot create pod")
		return
	}
	kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)

	time.Sleep(6 * time.Second)

	kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)

	_ = kubeletcontroller.StopPod(podConfig)
}

func TestMetricsTest(t *testing.T) {
	podConfig := GeneratePodConfigPy()
	err := kubeletcontroller.CreatePod(&podConfig)
	//err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot create pod")
		kubeletcontroller.StopPod(podConfig)
		return
	}
	_ = kubeletcontroller.NodeMetrics()
	//text := utils.JsonMarshal(stat)
	kubeletcontroller.StopPod(podConfig)
}
