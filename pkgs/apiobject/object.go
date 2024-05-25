package core

import "encoding/json"

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
	Controller bool    `yaml:"controller"`
}
type VersionLabel string

const (
	DELETE VersionLabel = "delete"
	UPDATE VersionLabel = "update"
	CREATE VersionLabel = "create"
)

// MarshalJSONList the object list to json
func MarshalJSONList(list interface{}) ([]byte, error) {
	return json.Marshal(list)
}
