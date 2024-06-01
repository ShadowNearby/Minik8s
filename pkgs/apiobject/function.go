package core

type Function struct {
	Kind       string       `yaml:"kind" json:"kind,omitempty"`
	APIVersion string       `yaml:"apiVersion" json:"apiVersion,omitempty"`
	Status     VersionLabel `yaml:"status" json:"status,omitempty"`
	/*name of functions*/
	Name string `yaml:"name" json:"name"`
	/*path to get functions*/
	Path string `yaml:"path" json:"path"`
}
