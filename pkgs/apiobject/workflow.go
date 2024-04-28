package core

type WorkflowSpec struct {
	EntryParams   string         `json:"entryParams" yaml:"entryParams"`
	EntryNodeName string         `json:"entryNodeName" yaml:"entryNodeName"`
	WorkflowNodes []WorkflowNode `json:"workflowNodes" yaml:"workflowNodes"`
}

type Workflow struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      WorkflowSpec `json:"spec" yaml:"spec"`
}

type WorkflowNode struct {
	Name       string             `json:"name" yaml:"name"`
	Type       WorkflowNodeType   `json:"type" yaml:"type"`
	FuncData   WorkflowFuncData   `json:"funcData" yaml:"funcData"`
	ChoiceData WorkflowChoiceData `json:"choiceData" yaml:"choiceData"`
}

type WorkflowChoiceData struct {
	TrueNextNodeName  string `json:"trueNextNodeName" yaml:"trueNextNodeName"`
	FalseNextNodeName string `json:"falseNextNodeName" yaml:"falseNextNodeName"`

	CheckType    ChoiceCheckType `json:"checkType" yaml:"checkType"`
	CheckVarName string          `json:"checkVarName" yaml:"checkVarName"`
	// 需要保证能够从上一个结果中获取到,填写json的key

	CompareValue string `json:"compareValue" yaml:"compareValue"` // 需要比较的值(无论是数字还是字符串，都需要转化为字符串)
}

type WorkflowNodeType string

const (
	WorkflowNodeTypeFunc   WorkflowNodeType = "func"
	WorkflowNodeTypeChoice WorkflowNodeType = "choice"

	WorkflowRunning   string = "running"
	WorkflowCompleted string = "completed"
)

type ChoiceCheckType string

const (
	ChoiceCheckTypeNumEqual               ChoiceCheckType = "numEqual"
	ChoiceCheckTypeNumNotEqual            ChoiceCheckType = "numNotEqual"
	ChoiceCheckTypeNumGreaterThan         ChoiceCheckType = "numGreaterThan"
	ChoiceCheckTypeNumLessThan            ChoiceCheckType = "numLessThan"
	ChoiceCheckTypeNumGreaterAndEqualThan ChoiceCheckType = "numGreaterAndEqualThan"
	ChoiceCheckTypeNumLessAndEqualThan    ChoiceCheckType = "numLessAndEqualThan"

	ChoiceCheckTypeStrEqual               ChoiceCheckType = "strEqual"
	ChoiceCheckTypeStrNotEqual            ChoiceCheckType = "strNotEqual"
	ChoiceCheckTypeStrGreaterThan         ChoiceCheckType = "strGreaterThan"
	ChoiceCheckTypeStrLessThan            ChoiceCheckType = "strLessThan"
	ChoiceCheckTypeStrGreaterAndEqualThan ChoiceCheckType = "strGreaterAndEqualThan"
	ChoiceCheckTypeStrLessAndEqualThan    ChoiceCheckType = "strLessAndEqualThan"
)

type WorkflowFuncData struct {
	FuncName      string `json:"funcName" yaml:"funcName"`
	FuncNamespace string `json:"funcNamespace" yaml:"funcNamespace"`
	NextNodeName  string `json:"nextNodeName" yaml:"nextNodeName"`
}

type WorkflowStatus struct {
	Phase  string `json:"phase" yaml:"phase"`
	Result string `json:"result" yaml:"result"`
}

type WorkflowStore struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      WorkflowSpec   `json:"spec" yaml:"spec"`
	Status    WorkflowStatus `json:"status" yaml:"status"`
}
