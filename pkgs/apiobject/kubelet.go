package core

import (
	"github.com/docker/go-connections/nat"
	"google.golang.org/grpc/resolver"
)

type Pod struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	MetaData   MetaData `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
	Status     Status   `json:"status" yaml:"status"`
}

type BasicInfo struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	MetaData   MetaData `json:"metadata" yaml:"metadata"`
}

type MetaData struct {
	Name            string            `json:"name" yaml:"name"`
	NameSpace       string            `json:"name_space" yaml:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	ResourceVersion string            `json:"resourceVersion" yaml:"resourceVersion,omitempty"`
	Annotations     map[string]string `json:"annotations"`
	UUID            string            `json:"uuid" yaml:"uuid"`
}

type Spec struct {
	Containers      []Container       `json:"containers" yaml:"containers"`
	RestartPolicy   restartPolicy     `json:"restartPolicy" yaml:"restartPolicy"`
	DnsPolicy       dnsPolicy         `json:"dnsPolicy,omitempty" yaml:"dnsPolicy,omitempty"`
	NodeSelector    map[string]string `json:"selector" yaml:"selector"`
	MinReadySeconds minReadySeconds   `json:"minReadySeconds,omitempty" yaml:"minReadySeconds,omitempty"`
	Selector        map[string]string `yaml:"selector"`
}

type Status struct {
	Phase          PhaseLabel     `json:"phase" yaml:"phase"`
	HostIP         string         `json:"host_ip" yaml:"hostIP"`
	PodIP          string         `json:"pod_ip" yaml:"podIP"`
	OwnerReference ownerReference `json:"owner_reference" yaml:"ownerReference"`
}
type restartPolicy string
type PhaseLabel string
type dnsPolicy string
type minReadySeconds int

type ownerReference struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Controller bool   `json:"controller,omitempty"`
}

type Container struct {
	Name            string              `json:"name" yaml:"name"`
	Image           string              `json:"image" yaml:"image"`
	ImagePullPolicy ImagePullPolicy     `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Cmd             []string            `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Args            []string            `json:"args,omitempty" yaml:"args,omitempty"`
	WorkingDir      string              `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	VolumeMounts    []VolumeMountConfig `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
	PortBindings    nat.PortMap         `json:"ports,omitempty" yaml:"ports,omitempty"` /* mapping of port bindings: container port -> []host ip+port */
	ExposedPorts    []string            `json:"exposedPorts" yaml:"exposedPorts"`       /* container's exposed ports */
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
	PodName        string
}

// const values

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

type ServiceStatus struct {
	Endpoints []resolver.Endpoint
	Phase     PhaseLabel `json:"phase"`
}

type ServicePort struct {
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	NodePort   int    `yaml:"node_port"`
	Protocol   string `yaml:"protocol"`
	TargetPort int    `yaml:"target_port"`
}

type ServiceSpec struct {
	Selector                      map[string]string `yaml:"selector"`
	Ports                         []ServicePort     `yaml:"ports"`
	AllocateLoadBalancerNodePorts bool              `yaml:"allocate_load_balancer_node_ports"`
	Type                          string            `yaml:"type"`
	ClusterIP                     string            `yaml:"cluster_ip"`
	ClusterIPs                    []string          `yaml:"cluster_ips"`
}

type Service struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      ServiceSpec `json:"spec" yaml:"spec"`
}

type ServiceStore struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      ServiceSpec   `json:"spec" yaml:"spec"`
	Status    ServiceStatus `json:"status" yaml:"status"`
}

type Node struct {
	ApiVersion   string       `json:"api_version,omitempty"`
	Kind         string       `json:"kind,omitempty"`
	NodeMetaData NodeMetaData `json:"metadata,omitempty"`
	Spec         NodeSpec     `json:"spec,omitempty"`
}

type NodeMetaData struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type NodeSpec struct {
	PodCIDR  string   `json:"podCIDR,omitempty"`
	PodCIDRs []string `json:"podCIDRs,omitempty"`
	Taints   []Taint  `json:"taints,omitempty"`
}

type Taint struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}

type NodeMetrics struct {
	Ready              bool   `json:"ready"`
	CPUUsage           uint64 `json:"cpuUsage"`
	MemoryUsage        uint64 `json:"memoryUsage"`
	PIDUsage           uint64 `json:"PIDUsage"`
	DiskUsage          uint64 `json:"diskUsage"`
	NetworkUnavailable bool   `json:"networkUnavailable"`
}

type KubeletConfig struct {
	MasterIP   string            `json:"masterIP"`
	MasterPort string            `json:"masterPort"`
	Labels     map[string]string `json:"labels"`
}
