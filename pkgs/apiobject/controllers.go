package core

import "encoding/json"

type ServiceStatus struct {
	Phase PhaseLabel `json:"phase"`
}

type ServicePort struct {
	Name       string `json:"name" yaml:"name"`
	Port       uint32 `json:"port" yaml:"port"`
	NodePort   uint32 `json:"nodePort" yaml:"nodePort"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	TargetPort string `json:"targetPort" yaml:"targetPort"`
}

type ServiceSpec struct {
	Selector  Selector      `json:"selector" yaml:"selector"`
	Ports     []ServicePort `json:"ports" yaml:"ports"`
	Type      ServiceType   `json:"type" yaml:"type"`
	ClusterIP string        `json:"clusterIP" yaml:"clusterIP"`
}

type Service struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      ServiceSpec   `json:"spec" yaml:"spec"`
	Status    ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ServiceType string

const (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the portal IP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"
)

func (p Service) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

/*---------------------------------ReplicaSet Types--------------------------------*/

type ReplicaSet struct {
	ApiVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind,omitempty" json:"kind"`
	MetaData   MetaData         `yaml:"metadata,omitempty" json:"metadata"`
	Spec       ReplicaSetSpec   `yaml:"spec,omitempty" json:"spec"`
	Status     ReplicaSetStatus `yaml:"status,omitempty" json:"status"`
}

type ReplicaSetSpec struct {
	Replicas int                `yaml:"replicas" json:"replicas"`
	Selector Selector           `yaml:"selector" json:"selector"`
	Template ReplicaSetTemplate `yaml:"template" json:"template"`
}

type ReplicaSetTemplate struct {
	MetaData MetaData `yaml:"metadata" json:"metadata"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
}

type ReplicaSetStatus struct {
	RealReplicas   int            `json:"realReplicas" yaml:"realReplicas"`
	Scale          int32          `json:"scale,omitempty"`
	OwnerReference OwnerReference `json:"ownerReference,omitempty"`
	// the truly ready replicas.
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

////TODO:Real replicas是指什么？
