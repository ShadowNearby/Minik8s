package core

type Endpoint struct {
	MetaData MetaData         `json:"metadata" yaml:"metadata"`
	Subsets  []EndpointSubset `json:"subsets" yaml:"subsets"`
}

type EndpointSubset struct {
	Addresses []EndpointAddress `json:"addresses" yaml:"addresses"`
	Ports     []EndpointPort    `json:"ports" yaml:"ports"`
}

type EndpointAddress struct {
	IP       string `json:"ip" yaml:"ip"`
	Hostname string `json:"hostname" yaml:"hostname"`
	NodeName string `json:"nodeName" yaml:"nodeName"`
}

type EndpointPort struct {
	Port     uint32 `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
	Name     string `json:"name" yaml:"name"`
}
