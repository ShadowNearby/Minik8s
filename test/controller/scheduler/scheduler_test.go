package scheduler

import (
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	scheduler "minik8s/pkgs/controller/scheduler"
	test "minik8s/test/kubelet"
	"testing"
)

func TestSchedule(t *testing.T) {
	var mockPod = test.GeneratePodConfigPy()
	mockPod.Spec.Selector = core.Selector{MatchLabels: map[string]string{"test": "haha"}}
	var myScheduler = scheduler.Scheduler{
		Policy: config.PolicyCPU,
	}
	ip, err := myScheduler.Schedule(mockPod)
	if err != nil {
		t.Error("schedule failed")
	}
	if ip != "172.31.184.139" && ip != "127.0.0.1" { // TODO
		t.Errorf("wrong ip: %s", ip)
	}
	mockPod.Spec.Selector = core.Selector{MatchLabels: map[string]string{"test": "nothaha"}}
	ip, err = myScheduler.Schedule(mockPod)
	if err == nil {
		t.Error("schedule failed")
	}
	if ip != "" {
		t.Error("expect nil")
	}
}
