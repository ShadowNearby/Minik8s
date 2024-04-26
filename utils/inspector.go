package utils

import (
	core "minik8s/pkgs/apiobject"
	"time"
)

func InitPodStatus(podConfig *core.Pod, nodeName string, podScheduled bool) core.PodStatus {
	// init a new pod in kubelet
	// construct containers status
	var ctStats = make(map[string]core.ContainerStatus, len(podConfig.Spec.Containers))
	for _, ctConfig := range podConfig.Spec.Containers {
		ctStats[GenerateContainerIDByName(ctConfig.Name, podConfig.MetaData.UUID)] = core.ContainerStatus{
			ContainerID: GenerateContainerIDByName(ctConfig.Name, podConfig.MetaData.UUID),
			Image:       ctConfig.Image,
			Port: core.PortStatus{
				PortNum:  0,
				Protocol: "TCP",
			}, // TODO: port status unknown
			State:        core.State{State: core.PhasePending},
			LastState:    core.State{},
			Ready:        false,
			RestartCount: 0,
			Environment:  "",
			Mounts:       ctConfig.VolumeMounts,
		}
	}
	var podStat = core.PodStatus{
		Name:       podConfig.MetaData.Name,
		Namespace:  podConfig.MetaData.NameSpace,
		Node:       nodeName,
		StartTime:  time.Now(),
		Labels:     podConfig.MetaData.Labels,
		Status:     core.PhasePending,
		Containers: ctStats,
		Conditions: core.Conditions{
			Initialized:     true,
			Ready:           false,
			ContainersReady: false,
			PodScheduled:    podScheduled,
		},
		NodeSelector: podConfig.Spec.NodeSelector,
	}
	return podStat
}
