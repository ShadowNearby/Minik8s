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
	ObjJob        ObjType = "job"
	ObjHpa        ObjType = "hpa"
	ObjFunction   ObjType = "function"
	ObjWorkflow   ObjType = "workflow"
	ObjDeployment ObjType = "deployment"
	ObjEndPoint   ObjType = "endpoint"
)

var ObjTypeAll = []string{
	"pod",
	"node",
	"replicaset",
	"service",
	"job",
	"hpa",
	"function",
	"workflow"}

var ObjTypeToCoreObjMap = map[ObjType]reflect.Type{
	ObjPod:        reflect.TypeOf(&Pod{}).Elem(),
	ObjNode:       reflect.TypeOf(&Node{}).Elem(),
	ObjReplicaSet: reflect.TypeOf(&ReplicaSet{}).Elem(),
	ObjService:    reflect.TypeOf(&Service{}).Elem(),
}

type ApiObjectKind interface {
	GetNameSpace() string
	GetObjectNamespace() string
}

func (p *Pod) GetNameSpace() string {
	return p.MetaData.Namespace
}
func (n *Node) GetNameSpace() string {
	return n.NodeMetaData.Namespace
}
func (w *WorkflowStore) GetNameSpace() string {
	return w.MetaData.Namespace
}
func (f *ReplicaSet) GetNameSpace() string {
	return f.MetaData.Namespace
}
