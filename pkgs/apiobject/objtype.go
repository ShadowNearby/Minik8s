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
	ObjEndPoint   ObjType = "endpoints"
	ObjDNS        ObjType = "dns"
	ObjVolume     ObjType = "volumes"
	ObjCsiVolume  ObjType = "csivolumes"
	ObjTask       ObjType = "tasks"    // ticktok, PingSource
	ObjTrigger    ObjType = "triggers" // only request
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
	"tasks",
	"triggers",
}

var ObjTypeToCoreObjMap = map[ObjType]reflect.Type{
	ObjPod:        reflect.TypeOf(&Pod{}).Elem(),
	ObjNode:       reflect.TypeOf(&Node{}).Elem(),
	ObjReplicaSet: reflect.TypeOf(&ReplicaSet{}).Elem(),
	ObjService:    reflect.TypeOf(&Service{}).Elem(),
	ObjHpa:        reflect.TypeOf(&HorizontalPodAutoscaler{}).Elem(),
	ObjFunction:   reflect.TypeOf(&Function{}).Elem(),
	ObjWorkflow:   reflect.TypeOf(&Workflow{}).Elem(),
	ObjEndPoint:   reflect.TypeOf(&Endpoint{}).Elem(),
	ObjDNS:        reflect.TypeOf(&DNSRecord{}).Elem(),
	ObjVolume:     reflect.TypeOf(&PersistentVolume{}).Elem(),
	ObjTrigger:    reflect.TypeOf(&TriggerMessage{}).Elem(),
	ObjTask:       reflect.TypeOf(&PingSource{}).Elem(),
	ObjJob:        reflect.TypeOf(&Job{}).Elem(),
}

var ObjTypeNamespace = map[ObjType]bool{
	ObjNode:       false,
	ObjFunction:   false,
	ObjWorkflow:   false,
	ObjVolume:     false,
	ObjCsiVolume:  false,
	ObjTask:       false,
	ObjTrigger:    false,
	ObjJob:        true,
	ObjPod:        true,
	ObjReplicaSet: true,
	ObjService:    true,
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
func (r *ReplicaSet) GetNamespace() string {
	return r.MetaData.Namespace
}
func (s *HorizontalPodAutoscaler) GetNamespace() string {
	return s.MetaData.Namespace
}
func (s *Endpoint) GetNamespace() string {
	return s.MetaData.Namespace
}
func (s *DNSRecord) GetNamespace() string {
	return s.MetaData.Namespace
}
func (s *Service) GetNamespace() string {
	return s.MetaData.Namespace
}
func (s *Job) GetNamespace() string { return s.MetaData.Namespace }
