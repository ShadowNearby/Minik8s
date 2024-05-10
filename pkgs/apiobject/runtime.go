package core

import (
	"time"
)

// CPUSet is a thread-safe, immutable set-like data structure for CPU IDs.
type CPUSet struct {
	elems map[int]struct{}
}

// CPUInfo contains the NUMA, socket, and core IDs associated with a CPU.
type CPUInfo struct {
	NUMANodeID int
	SocketID   int
	CoreID     int
}

// CPUDetails is a map from CPU ID to Core ID, Socket ID, and NUMA ID.
type CPUDetails map[int]CPUInfo

// CPUTopology contains details of node cpu, where :
// CPU  - logical CPU, cadvisor - thread
// Core - physical CPU, cadvisor - Core
// Socket - socket, cadvisor - Socket
// NUMA Node - NUMA cell, cadvisor - Node
type CPUTopology struct {
	NumCPUs    int
	NumCores   int
	NumSockets int
	CPUDetails CPUDetails
}

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
