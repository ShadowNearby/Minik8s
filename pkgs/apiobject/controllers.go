package core

import "google.golang.org/grpc/resolver"

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
	ClusterIP                     string            `yaml:"clusterIp"`
	ClusterIPs                    []string          `yaml:"clusterIps"`
}

type Service struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      ServiceSpec   `json:"spec" yaml:"spec"`
	Status    ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

/*---------------------------------ReplicaSet Types--------------------------------*/

type ReplicaSet struct {
	ApiVersion string         `yaml:"apiVersion" json:"apiVersion"`
	Kind       string         `yaml:"kind" json:"kind"`
	MetaData   MetaData       `yaml:"metaData" json:"metaData"`
	Spec       ReplicaSetSpec `yaml:"spec" json:"spec"`
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
	Replicas     int `json:"replicas" yaml:"replicas"`
	RealReplicas int `json:"realReplicas" yaml:"realReplicas"`
}
