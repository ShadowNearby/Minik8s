package core

type WorkState interface {
}

// for parsing
type RawState struct {
	Type StateType `json:"type"`
}

type TaskState struct {
	Type       StateType `json:"type" yaml:"type"`
	InputPath  string    `json:"inputPath,omitempty" yaml:"inputPath,omitempty"`
	ResultPath string    `json:"outputPath,omitempty" yaml:"outputPath,omitempty"`
	Next       string    `json:"next,omitempty" yaml:"next,omitempty"`
	End        bool      `json:"end,omitempty" yaml:"end,omitempty"`
}

type FailState struct {
	Type  StateType `json:"type" yaml:"type"`
	Error string    `json:"error" yaml:"error"`
	Cause string    `json:"cause" yaml:"cause"`
}

type ChoiceState struct {
	Type    StateType    `json:"type" yaml:"type"`
	Choices []ChoiceItem `json:"choices" yaml:"choices"`
	Default string       `json:"default,omitempty" yaml:"default,omitempty"`
}

type ChoiceItem struct {
	Variable                  string `json:"variable" yaml:"variable"`
	NumericEquals             *int   `json:"NumericEqual,omitempty" yaml:"NumericEqual,omitempty"`
	NumericNotEquals          *int   `json:"NumericNotEqual,omitempty" yaml:"NumericNotEqual,omitempty"`
	NumericLessThan           *int   `json:"NumericLessThan,omitempty" yaml:"NumericLessThan,omitempty"`
	NumericGreaterThan        *int   `json:"NumericGreaterThan,omitempty" yaml:"NumericGreaterThan,omitempty"`
	NumericLessThanOrEqual    *int   `json:"NumericLessThanOrEqual,omitempty" yaml:"NumericLessThanOrEqual,omitempty"`
	NumericGreaterThanOrEqual *int   `json:"NumericGreaterThanOrEqual,omitempty" yaml:"NumericGreaterThanOrEqual,omitempty"`

	StringEquals             *string `json:"StringEquals,omitempty" yaml:"StringEquals,omitempty"`
	StringNotEquals          *string `json:"StringNotEquals,omitempty" yaml:"StringNotEquals,omitempty"`
	StringLessThan           *string `json:"StringLessThan,omitempty" yaml:"StringLessThan,omitempty"`
	StringGreaterThan        *string `json:"StringGreaterThan,omitempty" yaml:"StringGreaterThan,omitempty"`
	StringLessThanOrEqual    *string `json:"StringLessThanOrEqual,omitempty" yaml:"StringLessThanOrEqual,omitempty"`
	StringGreaterThanOrEqual *string `json:"StringGreaterThanOrEqual,omitempty" yaml:"StringGreaterThanOrEqual,omitempty"`
	Next                     string  `json:"next" yaml:"next"`
}

type StateType string

const (
	Task     StateType = "TaskState"
	Choice   StateType = "ChoiceState"
	Parallel StateType = "Parallel"
	Wait     StateType = "Wait"
	Fail     StateType = "FailState"
	Succeed  StateType = "Succeed"
)

type Workflow struct {
	APIVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind       string       `json:"kind" yaml:"kind"`
	Name       string       `json:"name" yaml:"name"`
	Status     VersionLabel `json:"status,omitempty" yaml:"status,omitempty"`
	StartAt    string       `json:"startAt" yaml:"startAt"`

	States map[string]WorkState `json:"states" yaml:"states"`

	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
}

type WorkFlowTriggerRequest struct {
	Name   string `json:"name" yaml:"name"`
	Params string `json:"params" yaml:"params"`
}
