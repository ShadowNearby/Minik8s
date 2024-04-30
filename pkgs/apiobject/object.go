package core

import "minik8s/pkgs/config"

type MetaData struct {
	Name            string            `json:"name" yaml:"name"`
	NameSpace       string            `json:"nameSpace" yaml:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	ResourceVersion string            `json:"resourceVersion" yaml:"resourceVersion,omitempty"`
	Annotations     map[string]string `json:"annotations"`
	UUID            string            `json:"uuid" yaml:"uuid"`
	OwnerReference  OwnerReference    `json:"ownerReference" yaml:"ownerReference"`
}
type Selector struct {
	MatchLabels map[string]string `yaml:"matchLabels" json:"matchLabels"`
}
type OwnerReference struct {
	ObjType    config.ObjType `json:"objType"`
	Name       string         `json:"name"`
	NameSpace  string         `json:"nameSpace"`
	Controller bool           `yaml:"controller"`
}
