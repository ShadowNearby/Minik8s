package core

import (
	"github.com/docker/go-connections/nat"
)

type Pod struct {
	ApiVersion string   `json:"api_version"`
	MetaData   MetaData `json:"meta_data"`
	Spec       Spec     `json:"Spec"`
	Status     Status   `json:"Status"`
}

type MetaData struct {
	Name            string            `json:"name"`
	NameSpace       string            `json:"name_space"`
	Labels          map[string]string `json:"labels,omitempty"`
	ResourceVersion string            `json:"resource_version"`
	UUID            string
}

type Spec struct {
	Containers    []Container       `json:"containers"`
	RestartPolicy restartPolicy     `json:"restart_policy"`
	NodeSelector  map[string]string `json:"node_selector"`
}

type Status struct {
	Phase          phaseLabel     `json:"phase"`
	HostIP         string         `json:"host_ip"`
	PodIP          string         `json:"pod_ip"`
	OwnerReference ownerReference `json:"owner_reference"`
}
type restartPolicy string
type phaseLabel string

type ownerReference struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Controller bool   `json:"controller,omitempty"`
}

type Container struct {
	Name            string              `json:"name"`
	Image           string              `json:"image"`
	ImagePullPolicy ImagePullPolicy     `json:"image_pull_policy,omitempty"`
	Cmd             []string            `json:"cmd,omitempty"`
	Args            []string            `json:"args,omitempty"`
	WorkingDir      string              `json:"working_dir,omitempty"`
	VolumeMounts    []VolumeMountConfig `json:"volume_mounts,omitempty"`
	PortBindings    nat.PortMap         `json:"ports,omitempty"` /* mapping of port bindings: container port -> []host ip+port */
	ExposedPorts    []string            `json:"exposed_ports"`   /* container's exposed ports */
	Env             []EnvConfig         `json:"env,omitempty"`
	Resources       ResourcesConfig     `json:"resources,omitempty"`
}

type ImagePullPolicy string

type VolumeMountConfig struct {
	ContainerPath string `json:"container_path"`
	HostPath      string `json:"mount_path"`
	ReadOnly      bool   `json:"read_only"`
}

type PortConfig struct {
	Name          string `json:"name"`
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
	Protocol      string `json:"protocol"`
}

type EnvConfig struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Limit struct {
	Cpu    string `json:"cpu"` /* 0-3, 0, 1 */
	Memory uint64 `json:"memory"`
}

type Request struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type ResourcesConfig struct {
	Limit   Limit   `json:"limit"`
	Request Request `json:"Request"`
}

type ContainerdSpec struct {
	Namespace      string
	Image          string
	Name           string
	ID             string
	VolumeMounts   map[string]string
	Cmd            []string
	Envs           []string
	Resource       Limit
	Labels         map[string]string
	LinuxNamespace map[string]string /* support localhost communication */
	PullPolicy     ImagePullPolicy
}

// Inspect inspect data structure
type Inspect struct {
	State           InspectState
	ResolveConfPath string
}

type InspectState struct {
	Status     phaseLabel
	Running    bool
	Paused     bool
	Restarting bool
	Pid        uint64
}

// const values

const (
	EmptyCpu    string = ""
	EmptyMemory uint64 = 0
)
const (
	PhasePending phaseLabel = "pending"
	PhaseRunning phaseLabel = "running"
	PhaseSucceed phaseLabel = "succeed"
	PhaseFailed  phaseLabel = "failed"
	PhaseUnknown phaseLabel = "unknown"
)

const (
	PullAlways  ImagePullPolicy = "Always"
	PullNever   ImagePullPolicy = "Never"
	PullIfNeeds ImagePullPolicy = "IfNotPresent"
)

const (
	RestartAlways    restartPolicy = "Always"
	RestartOnFailure restartPolicy = "OnFailure"
)
