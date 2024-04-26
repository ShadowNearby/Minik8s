package core

import (
	"reflect"
	"strings"
)

const (
	PodKind        = "Pod"
	ServiceKind    = "Service"
	DnsKind        = "Dns"
	NodeKind       = "Node"
	JobKind        = "Job"
	ReplicaSetKind = "Replicaset"
	HpaKind        = "Hpa"
	FunctionKind   = "Function"
	WorkflowKind   = "Workflow"
)

var KindToStructType = map[string]reflect.Type{
	PodKind:     reflect.TypeOf(&Pod{}).Elem(),
	ServiceKind: reflect.TypeOf(&Service{}).Elem(),
	//DnsKind:        reflect.TypeOf(&Dns{}).Elem(),
	//JobKind:        reflect.TypeOf(&Job{}).Elem(),
	//NodeKind:       reflect.TypeOf(&Node{}).Elem(),
	//ReplicaSetKind: reflect.TypeOf(&ReplicaSet{}).Elem(),
	//HpaKind:        reflect.TypeOf(&HPA{}).Elem(),
	//FunctionKind:   reflect.TypeOf(&Function{}).Elem(),
	//WorkflowKind:   reflect.TypeOf(&Workflow{}).Elem(),
}
var AllResourceKindSlice = []string{PodKind, ServiceKind, DnsKind, NodeKind, JobKind, ReplicaSetKind, HpaKind, FunctionKind, WorkflowKind}

var AllResourceKind = strings.ToLower("[" + PodKind + "/" + ServiceKind + "/" + DnsKind + "/" + NodeKind + "/" + JobKind +
	"/" + ReplicaSetKind + "/" + HpaKind + "/" + FunctionKind + "/" + WorkflowKind + "]")

type ApiObjectKind interface {
	GetObjectKind() string
	GetObjectName() string
	GetObjectNamespace() string
}

func (w *WorkflowStore) GetName() string {
	return w.MetaData.Name
}

func (w *WorkflowStore) GetNamespace() string {
	return w.MetaData.NameSpace
}

func (wf *Workflow) GetObjectKind() string {
	return wf.Kind
}

func (wf *Workflow) GetObjectName() string {
	return wf.MetaData.Name
}

func (wf *Workflow) GetObjectNamespace() string {
	return wf.MetaData.NameSpace
}

func (w *Workflow) ToWorkflowStore() *WorkflowStore {
	// 创建一个Status是空的WorkflowStore
	return &WorkflowStore{
		BasicInfo: w.BasicInfo,
		Spec:      w.Spec,
		Status:    WorkflowStatus{},
	}
}

// 定义WorkflowStore转化为Workflow的函数
func (w *WorkflowStore) ToWorkflow() *Workflow {
	// 创建一个Status是空的WorkflowStore
	return &Workflow{
		BasicInfo: w.BasicInfo,
		Spec:      w.Spec,
	}
}
