package utils

import (
	core "minik8s/pkgs/apiobject"
	"time"
)

func InitPodStatus(podConfig *core.Pod) core.PodStatus {
	var podStat = core.PodStatus{
		HostIP:           podConfig.Status.HostIP,
		PodIP:            "",
		StartTime:        time.Now(),
		OldStatus:        make([]core.Status, 0),
		ContainersStatus: make([]core.ContainerStatus, 0),
		Condition:        core.CondPending,
	}
	return podStat
}
