package core

import (
	"encoding/json"

	"github.com/docker/go-connections/nat"
)

type Pod struct {
	ApiVersion string    `json:"apiVersion" yaml:"apiVersion"`
	MetaData   MetaData  `json:"metadata" yaml:"metadata"`
	Spec       PodSpec   `json:"spec" yaml:"spec"`
	Status     PodStatus `json:"podStatus" yaml:"podStatus"`
}

func (p Pod) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

type BasicInfo struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	MetaData   MetaData `json:"metadata" yaml:"metadata"`
}

type PodSpec struct {
	Containers      []Container     `json:"containers" yaml:"containers"`
	RestartPolicy   restartPolicy   `json:"restartPolicy" yaml:"restartPolicy"`
	DnsPolicy       dnsPolicy       `json:"dnsPolicy,omitempty" yaml:"dnsPolicy,omitempty"`
	Selector        Selector        `json:"selector" yaml:"selector"`
	MinReadySeconds minReadySeconds `json:"minReadySeconds,omitempty" yaml:"minReadySeconds,omitempty"`
	Volumes         []Volume        `json:"volumes" yaml:"volumes"`
}

type restartPolicy string
type PhaseLabel string
type dnsPolicy string
type minReadySeconds int

//type OwnerReference struct {
//	ApiVersion string  `json:"apiVersion"`
//	Kind       ObjType `json:"kind"`
//	Name       string  `json:"name"`
//	UID        string  `json:"UID"`
//	Controller bool    `json:"controller,omitempty"` /* true means under control, use *bool? */
//}

type Container struct {
	Name            string              `json:"name" yaml:"name"`
	Image           string              `json:"image" yaml:"image"`
	ImagePullPolicy ImagePullPolicy     `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Cmd             []string            `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Args            []string            `json:"args,omitempty" yaml:"args,omitempty"`
	WorkingDir      string              `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	VolumeMounts    []VolumeMountConfig `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
	PortBindings    nat.PortMap         `json:"portsBindings,omitempty" yaml:"portsBindings,omitempty"` /* mapping of port bindings: container port -> []host ip+port */
	Ports           []PortConfig        `json:"ports,omitempty" yaml:"ports,omitempty"`
	ExposedPorts    []string            `json:"exposedPorts" yaml:"exposedPorts"` /* container's exposed ports */
	Env             []EnvConfig         `json:"env,omitempty" yaml:"env,omitempty"`
	Resources       ResourcesConfig     `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type ImagePullPolicy string

type VolumeMountConfig struct {
	ContainerPath string `json:"container_path"`
	HostPath      string `json:"mount_path"`
	ReadOnly      bool   `json:"read_only"`
}

type PortConfig struct {
	Name          string `json:"name" yaml:"name"`
	ContainerPort uint32 `json:"containerPort" yaml:"containerPort"`
	HostPort      string `json:"hostPort" yaml:"hostPort"`
	Protocol      string `json:"protocol" yaml:"protocol"`
}

type EnvConfig struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type ResourceLimit struct {
	Cpu     string `json:"cpu"` /* 0-3, 0, 1 */
	Memory  uint64 `json:"memory"`
	Storage uint64 `json:"storage"`
}

type ResourcesConfig struct {
	Limit   ResourceLimit `json:"limit"`
	Request ResourceLimit `json:"Request"`
}

type ContainerdSpec struct {
	Namespace      string
	Image          string
	Name           string
	ID             string
	VolumeMounts   map[string]string
	Cmd            []string
	Args           []string
	Envs           []string
	Resource       ResourceLimit
	Labels         map[string]string
	LinuxNamespace map[string]string /* support localhost communication */
	PullPolicy     ImagePullPolicy
	PodName        string
}

// Inspect inspect data structure
type Inspect struct {
	State           InspectState
	ResolveConfPath string
}

type InspectState struct {
	Status     PhaseLabel
	Running    bool
	Paused     bool
	Restarting bool
	Pid        uint64
}

type Node struct {
	ApiVersion   string   `json:"apiVersion,omitempty"`
	Kind         string   `json:"kind,omitempty"`
	NodeMetaData MetaData `json:"metadata,omitempty"`
	Spec         NodeSpec `json:"spec,omitempty"`
}

func (p Node) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

type NodeSpec struct {
	PodCIDR  string   `json:"podCIDR,omitempty"`
	PodCIDRs []string `json:"podCIDRs,omitempty"`
	NodeIP   string   `json:"nodeIP"`
	Taints   []Taint  `json:"taints,omitempty"`
}

type Taint struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}
type KubeletConfig struct {
	MasterIP   string            `json:"masterIP"`
	MasterPort string            `json:"masterPort"`
	Labels     map[string]string `json:"labels"`
}
