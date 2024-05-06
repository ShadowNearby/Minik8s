package core

type Endpoint struct {
	MetaData         MetaData `json:"metadata" yaml:"metadata"`
	ServiceClusterIP string   `json:"serviceClusterIP" yaml:"serviceClusterIP"`
	Binds            []EndpointBind
}

type EndpointBind struct {
	ServicePort  uint32                `json:"servicePort" yaml:"servicePort"`
	Destinations []EndpointDestination `json:"destinations" yaml:"destinations"`
}

type EndpointDestination struct {
	IP   string `json:"ip" yaml:"ip"`
	Port uint32 `json:"port" yaml:"port"`
}
