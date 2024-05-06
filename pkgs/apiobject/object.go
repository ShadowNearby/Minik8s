package core

type MetaData struct {
	Name            string            `json:"name" yaml:"name"`
	Namespace       string            `json:"namespace" yaml:"namespace,omitempty"`
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
	ObjType    ObjType `json:"objType"`
	Name       string  `json:"name"`
	Namespace  string  `json:"namespace"`
	Controller bool    `yaml:"controller"`
}
