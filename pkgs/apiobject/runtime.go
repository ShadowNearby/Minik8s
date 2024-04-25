package core

import (
	"time"
)

type ContainerMetrics struct {
	PidCount    uint64
	CpuUsage    uint64
	MemoryUsage uint64
	DiskUsage   uint64
}

var EmptyContainerMetrics = ContainerMetrics{
	PidCount:    0,
	CpuUsage:    0,
	MemoryUsage: 0,
}

// container status is described in containerd.Status

type PortStatus struct {
	PortNum  uint64
	Protocol string
}

type State struct {
	State    PhaseLabel `json:"state"`
	Reason   string     `json:"reason"`
	ExitCode uint32     `json:"exit_code"`
	Started  time.Time  `json:"started"`
	Finished time.Time  `json:"finished"`
}

type ContainerStatus struct {
	ContainerID  string              `json:"container_id"`
	Image        string              `json:"image"`
	Port         PortStatus          `json:"port"`
	State        State               `json:"state"`
	LastState    State               `json:"last_state,omitempty"`
	Ready        bool                `json:"ready"`
	RestartCount uint64              `json:"restart_count"`
	Environment  string              `json:"environment"`
	Mounts       []VolumeMountConfig `json:"mounts,omitempty"` // TODO: not sure
}

type Conditions struct {
	Initialized     bool `json:"initialized"` // i dont think we should implement this
	Ready           bool `json:"ready"`
	ContainersReady bool `json:"containers_ready"`
	PodScheduled    bool `json:"pod_scheduled"`
}

type PodStatus struct {
	Name         string                     `json:"name"`
	Namespace    string                     `json:"namespace"`
	Node         string                     `json:"node"`
	StartTime    time.Time                  `json:"start_time"`
	Labels       map[string]string          `json:"labels"`
	Status       PhaseLabel                 `json:"phase"`
	Containers   map[string]ContainerStatus `json:"containers"`
	Conditions   Conditions                 `json:"conditions"`
	NodeSelector map[string]string          `json:"node_selector"`
	//Volumes map[string]string // TODO: dont know the type, should we implement them?
	//QosClass string
	//Tolerations []string
	//Events []string
}
