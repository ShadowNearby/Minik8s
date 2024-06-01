package core

import (
	"reflect"
)

type ObjType string

const (
	ObjPod        ObjType = "pods"
	ObjNode       ObjType = "nodes"
	ObjReplicaSet ObjType = "replicas"
	ObjService    ObjType = "services"
	ObjJob        ObjType = "jobs"
	ObjHpa        ObjType = "hpa"
	ObjFunction   ObjType = "functions"
	ObjWorkflow   ObjType = "workflows"
	ObjDeployment ObjType = "deployment"
	ObjEndPoint   ObjType = "endpoints"
	ObjDNS        ObjType = "dns"
	ObjVolume     ObjType = "volumes"
	ObjCsiVolume  ObjType = "csivolumes"
)

var ObjTypeAll = []string{
	"pods",
	"nodes",
	"replicas",
	"services",
	"endpoints",
	"deployment",
	"jobs",
	"hpa",
	"functions",
	"workflows",
	"dns",
	"volumes",
}

var ObjTypeToCoreObjMap = map[ObjType]reflect.Type{
	ObjPod:        reflect.TypeOf(&Pod{}).Elem(),
	ObjNode:       reflect.TypeOf(&Node{}).Elem(),
	ObjReplicaSet: reflect.TypeOf(&ReplicaSet{}).Elem(),
	ObjService:    reflect.TypeOf(&Service{}).Elem(),
	ObjDeployment: reflect.TypeOf(&ReplicaSet{}).Elem(),
	ObjEndPoint:   reflect.TypeOf(&Endpoint{}).Elem(),
	ObjDNS:        reflect.TypeOf(&DNSRecord{}).Elem(),
	ObjFunction:   reflect.TypeOf(&Function{}).Elem(),
}

type ApiObjectKind interface {
	GetNameSpace() string
}

func (p *Pod) GetNameSpace() string {
	return p.MetaData.Namespace
}
func (n *Node) GetNameSpace() string {
	return n.NodeMetaData.Namespace
}
func (r *ReplicaSet) GetNameSpace() string {
	return r.MetaData.Namespace
}
func (s *Service) GetNameSpace() string {
	return s.MetaData.Namespace
}
