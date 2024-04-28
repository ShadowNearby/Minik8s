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
	NameSpace       string            `json:"nameSpace" yaml:"namespace,omitempty"`
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
	HostIP         string         `json:"hostIP" yaml:"hostIP"`
	PodIP          string         `json:"podIP" yaml:"podIP"`
	OwnerReference ownerReference `json:"ownerReference" yaml:"ownerReference"`
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
	ContainerPath string `json:"containerPath" yaml:"containerPath"`
	HostPath      string `json:"mountPath" yaml:"mountPath"`
	ReadOnly      bool   `json:"readOnly" yaml:"readOnly"`
}

type PortConfig struct {
	Name          string `json:"name" yaml:"name"`
	ContainerPort string `json:"containerPort" yaml:"containerPort"`
	HostPort      string `json:"hostPort" yaml:"hostPort"`
	Protocol      string `json:"protocol" yaml:"protocol"`
}

type EnvConfig struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
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
	NodePort   int    `yaml:"nodePort"`
	Protocol   string `yaml:"protocol"`
	TargetPort int    `yaml:"target_port"`
}

type ServiceSpec struct {
	Selector                      map[string]string `yaml:"selector"`
	Ports                         []ServicePort     `yaml:"ports"`
	AllocateLoadBalancerNodePorts bool              `yaml:"allocateLoadBalancerNodePorts"`
	Type                          string            `yaml:"type"`
	ClusterIP                     string            `yaml:"clusterIp"`
	ClusterIPs                    []string          `yaml:"clusterIps"`
}

type Service struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      ServiceSpec   `json:"spec" yaml:"spec"`
	Status    ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type Node struct {
	ApiVersion   string       `json:"apiVersion,omitempty"`
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
type AutoScale struct {
	Name        string `yaml:"name"`
	Workload    Pod    `yaml:"workload"`
	MinReplicas int    `yaml:"minReplicas"`
	MaxReplicas int    `yaml:"maxReplicas"`
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
type WorkflowSpec struct {
	EntryParams   string         `json:"entryParams" yaml:"entryParams"`
	EntryNodeName string         `json:"entryNodeName" yaml:"entryNodeName"`
	WorkflowNodes []WorkflowNode `json:"workflowNodes" yaml:"workflowNodes"`
}

type Workflow struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      WorkflowSpec `json:"spec" yaml:"spec"`
}

type WorkflowNode struct {
	Name       string             `json:"name" yaml:"name"`
	Type       WorkflowNodeType   `json:"type" yaml:"type"`
	FuncData   WorkflowFuncData   `json:"funcData" yaml:"funcData"`
	ChoiceData WorkflowChoiceData `json:"choiceData" yaml:"choiceData"`
}

type WorkflowChoiceData struct {
	TrueNextNodeName  string `json:"trueNextNodeName" yaml:"trueNextNodeName"`
	FalseNextNodeName string `json:"falseNextNodeName" yaml:"falseNextNodeName"`

	CheckType    ChoiceCheckType `json:"checkType" yaml:"checkType"`
	CheckVarName string          `json:"checkVarName" yaml:"checkVarName"`
	// 需要保证能够从上一个结果中获取到,填写json的key

	CompareValue string `json:"compareValue" yaml:"compareValue"` // 需要比较的值(无论是数字还是字符串，都需要转化为字符串)
}

type WorkflowNodeType string

const (
	WorkflowNodeTypeFunc   WorkflowNodeType = "func"
	WorkflowNodeTypeChoice WorkflowNodeType = "choice"

	WorkflowRunning   string = "running"
	WorkflowCompleted string = "completed"
)

type ChoiceCheckType string

const (
	ChoiceCheckTypeNumEqual               ChoiceCheckType = "numEqual"
	ChoiceCheckTypeNumNotEqual            ChoiceCheckType = "numNotEqual"
	ChoiceCheckTypeNumGreaterThan         ChoiceCheckType = "numGreaterThan"
	ChoiceCheckTypeNumLessThan            ChoiceCheckType = "numLessThan"
	ChoiceCheckTypeNumGreaterAndEqualThan ChoiceCheckType = "numGreaterAndEqualThan"
	ChoiceCheckTypeNumLessAndEqualThan    ChoiceCheckType = "numLessAndEqualThan"

	ChoiceCheckTypeStrEqual               ChoiceCheckType = "strEqual"
	ChoiceCheckTypeStrNotEqual            ChoiceCheckType = "strNotEqual"
	ChoiceCheckTypeStrGreaterThan         ChoiceCheckType = "strGreaterThan"
	ChoiceCheckTypeStrLessThan            ChoiceCheckType = "strLessThan"
	ChoiceCheckTypeStrGreaterAndEqualThan ChoiceCheckType = "strGreaterAndEqualThan"
	ChoiceCheckTypeStrLessAndEqualThan    ChoiceCheckType = "strLessAndEqualThan"
)

type WorkflowFuncData struct {
	FuncName      string `json:"funcName" yaml:"funcName"`
	FuncNamespace string `json:"funcNamespace" yaml:"funcNamespace"`
	NextNodeName  string `json:"nextNodeName" yaml:"nextNodeName"`
}

type WorkflowStatus struct {
	Phase  string `json:"phase" yaml:"phase"`
	Result string `json:"result" yaml:"result"`
}

type WorkflowStore struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      WorkflowSpec   `json:"spec" yaml:"spec"`
	Status    WorkflowStatus `json:"status" yaml:"status"`
}
