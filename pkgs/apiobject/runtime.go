package core

import (
	"time"
)

type NodeMetrics struct {
	Ready              bool   `json:"ready"`
	CPUUsage           uint64 `json:"cpuUsage"`
	MemoryUsage        uint64 `json:"memoryUsage"`
	PIDUsage           uint64 `json:"PIDUsage"`
	DiskUsage          uint64 `json:"diskUsage"`
	NetworkUnavailable bool   `json:"networkUnavailable"`
}

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
	DiskUsage:   0,
}

type PodStatus struct {
	Phase            PhaseLabel        `json:"phase" yaml:"phase"`
	HostIP           string            `json:"hostIP" yaml:"hostIP"` /* node name */
	PodIP            string            `json:"podIP" yaml:"podIP"`
	StartTime        time.Time         `yaml:"startTime"`
	OldStatus        []Status          `json:"oldStatus"`
	ContainersStatus []ContainerStatus `json:"containersStatus"`
	Condition        Condition         `json:"condition"`
}

type ContainerStatus struct {
	ID           string         `json:"ID"`
	Name         string         `json:"name"`
	Image        string         `json:"image"`
	State        ContainerState `json:"state"`
	RestartCount int32          `json:"restartCount"`
	Environment  []EnvConfig    `json:"environment"`
	//Mounts       []VolumeMountConfig `json:"mounts,omitempty"`
}

type ContainerState string

const (
	ContainerWaiting    ContainerState = "waiting"
	ContainerRunning    ContainerState = "running"
	ContainerTerminated ContainerState = "terminated"
)

type Status struct {
	Reason   string    `json:"reason"`
	ExitCode uint32    `json:"exit_code"`
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`
}

type Condition string

const (
	ConScheduled       Condition = "pod scheduled"
	ConInitialized     Condition = "pod initialized"
	ConReady           Condition = "pod ready"
	ConContainersReady Condition = "pod containers ready"
)

const PauseContainerName string = "pause_container"
const (
	EmptyCpu    string = ""
	EmptyMemory uint64 = 0
)
const (
	PhasePending PhaseLabel = "pending"
	PhaseRunning PhaseLabel = "running"
	PhaseSucceed PhaseLabel = "succeed"
	PhaseFailed  PhaseLabel = "failed"
	PhaseUnknown PhaseLabel = "unknown"
)

const (
	PullAlways  ImagePullPolicy = "Always"
	PullNever   ImagePullPolicy = "Never"
	PullIfNeeds ImagePullPolicy = "IfNotPresent"
)

const (
	RestartAlways    restartPolicy = "Always"
	RestartOnFailure restartPolicy = "OnFailure"
	RestartNever     restartPolicy = "Never"
)
