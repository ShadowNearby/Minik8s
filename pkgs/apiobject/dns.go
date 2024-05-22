package core

type DNSRecord struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Host      string    `json:"host" yaml:"host"`
	Paths     []DNSPath `json:"paths" yaml:"paths"`
}

type DNSPath struct {
	IP      string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Port    uint32 `json:"port" yaml:"port"`
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
	Service string `json:"service" yaml:"service"`
}
type DNSEntry struct {
	Host string `json:"host" yaml:"host"`
}

type NginxLocation struct {
	Path string
	IP   string
	Port uint32
}
type NginxServer struct {
	Addr       string
	ServerName string
	Locations  []NginxLocation
}
type NginxConf struct {
	Servers []NginxServer
}
