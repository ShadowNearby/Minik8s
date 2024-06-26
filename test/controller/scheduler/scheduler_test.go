package scheduler

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	scheduler "minik8s/pkgs/controller/scheduler"
	"minik8s/utils"
	"testing"
)

func TestSchedule(t *testing.T) {
	var mockPod = utils.GeneratePodConfigPy()
	mockPod.Spec.Selector = core.Selector{MatchLabels: map[string]string{"test": "haha"}}
	var myScheduler = scheduler.Scheduler{
		Policy: constants.PolicyCPU,
	}
	ip, err := myScheduler.Schedule(mockPod)
	if err != nil {
		t.Error("schedule failed")
	}
	if ip != "192.168.1.12" && ip != "127.0.0.1" { // TODO
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
