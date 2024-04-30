package utils

import (
	core "minik8s/pkgs/apiobject"
	"os"
	"time"
)

func InitPodStatus(podConfig *core.Pod) core.PodStatus {
	host, _ := os.Hostname()
	var podStat = core.PodStatus{
		Phase:            core.PhasePending,
		HostIP:           host,
		PodIP:            "",
		StartTime:        time.Now(),
		OldStatus:        make([]core.Status, 0),
		ContainersStatus: make([]core.ContainerStatus, 0),
		Condition:        core.ConInitialized,
	}
	return podStat
}
