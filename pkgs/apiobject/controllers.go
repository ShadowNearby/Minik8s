package core

type ServiceStatus struct {
	Phase PhaseLabel `json:"phase"`
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
