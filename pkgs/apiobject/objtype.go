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
}

var ObjTypeNamespace = map[ObjType]bool{
	ObjNode:      false,
	ObjFunction:  false,
	ObjWorkflow:  false,
	ObjVolume:    false,
	ObjCsiVolume: false,

	ObjPod:        true,
	ObjReplicaSet: true,
	ObjService:    true,
	ObjJob:        true,
	ObjHpa:        true,
	ObjEndPoint:   true,
	ObjDNS:        true,
}

// {ObjNode, ObjFunction, ObjWorkflow, ObjVolume, ObjCsiVolume}

type ApiObjectKind interface {
	GetNamespace() string
}

func (p *Pod) GetNamespace() string {
	return p.MetaData.Namespace
}
func (n *Node) GetNamespace() string {
	return n.NodeMetaData.Namespace
}
func (w *Workflow) GetNamespace() string {
	return "workflow"
}
func (r *ReplicaSet) GetNamespace() string {
	return r.MetaData.Namespace
}
func (s *Service) GetNamespace() string {
	return s.MetaData.Namespace
}
