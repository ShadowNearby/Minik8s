package core

type ServiceStatus struct {
	Phase PhaseLabel `json:"phase"`
}

type ServicePort struct {
	Name       string `yaml:"name"`
	Port       uint32 `yaml:"port"`
	NodePort   uint32 `yaml:"nodePort"`
	Protocol   string `yaml:"protocol"`
	TargetPort string `yaml:"targetPort"`
}

type ServiceSpec struct {
	Selector                      map[string]string `yaml:"selector"`
	Ports                         []ServicePort     `yaml:"ports"`
	AllocateLoadBalancerNodePorts bool              `yaml:"allocateLoadBalancerNodePorts"`
	Type                          string            `yaml:"type"`
	ClusterIP                     string            `yaml:"clusterIp"`
}

type Service struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      ServiceSpec   `json:"spec" yaml:"spec"`
	Status    ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

/*---------------------------------ReplicaSet Types--------------------------------*/

type ReplicaSet struct {
	ApiVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind" json:"kind"`
	MetaData   MetaData         `yaml:"metaData" json:"metaData"`
	Spec       ReplicaSetSpec   `yaml:"spec" json:"spec"`
	Status     ReplicaSetStatus `yaml:"status" json:"status"`
}

type ReplicaSetSpec struct {
	Replicas int                `yaml:"replicas" json:"replicas"`
	Selector Selector           `yaml:"selector" json:"selector"`
	Template ReplicaSetTemplate `yaml:"template" json:"template"`
}

type ReplicaSetTemplate struct {
	MetaData MetaData `yaml:"metaData" json:"metaData"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
}

type ReplicaSetStatus struct {
	RealReplicas int `json:"realReplicas" yaml:"realReplicas"`
}
