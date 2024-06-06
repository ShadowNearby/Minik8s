package core

import (
	"time"
)

// NodeMetrics usage means percentage
type NodeMetrics struct {
	Ready              bool    `json:"ready"`
	CPUUsage           float64 `json:"cpuUsage"`
	MemoryUsage        float64 `json:"memoryUsage"`
	PIDUsage           float64 `json:"PIDUsage"`
	DiskUsage          float64 `json:"diskUsage"`
	NetworkUnavailable bool    `json:"networkUnavailable"`
}

type ContainerMetrics struct {
	CpuUsage    float64 // percentage
	MemoryUsage float64 // percentage
}

var EmptyContainerMetrics = ContainerMetrics{
	CpuUsage:    0,
	MemoryUsage: 0,
}

type PodStatus struct {
	HostIP           string            `json:"hostIP" yaml:"hostIP"` /* node ip */
	PodIP            string            `json:"podIP" yaml:"podIP"`
	StartTime        time.Time         `yaml:"startTime"`
	OldStatus        []Status          `json:"oldStatus"`
	ContainersStatus []ContainerStatus `json:"containersStatus"`
	Phase            PodPhase          `json:"phase"`
}

type ContainerStatus struct {
	ID             string                  `json:"ID"`
	Name           string                  `json:"name"`
	Image          string                  `json:"image"`
	State          ContainerState          `json:"state"`
	RestartCount   int32                   `json:"restartCount"`
	Environment    []EnvConfig             `json:"environment"`
	TerminateState ContainerTerminateState `json:"terminatedState"`
	//Mounts       []VolumeMountConfig `json:"mounts,omitempty"`
}

type ContainerState string

const (
	ContainerWaiting    ContainerState = "Waiting"
	ContainerRunning    ContainerState = "Running"
	ContainerTerminated ContainerState = "Terminated"
)

type ContainerTerminateState string

const (
	ContainerTerminateSuccess ContainerTerminateState = "success"
	ContainerTerminateFail    ContainerTerminateState = "fail"
)

type Status struct {
	Reason   string    `json:"reason"`
	ExitCode uint32    `json:"exitcode"`
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`
}

type PodPhase string

const (
	PodPhasePending   PodPhase = "Pending"
	PodPhaseRunning   PodPhase = "Running"
	PodPhaseSucceeded PodPhase = "Succeeded"
	PodPhaseFailed    PodPhase = "Failed"
	PodPhaseUnknown   PodPhase = "Unknown"
)

const PauseContainerName string = "pause-container"
const (
	EmptyCpu    string = ""
	EmptyMemory uint64 = 0
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
