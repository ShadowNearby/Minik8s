package core

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type WorkState interface {
}

type TaskState struct {
	Type       StateType `json:"type"`
	InputPath  string    `json:"inputPath,omitempty"`
	ResultPath string    `json:"outputPath,omitempty"`
	Next       string    `json:"next,omitempty"`
	End        bool      `json:"end,omitempty"`
}

type FailState struct {
	Type  StateType `json:"type"`
	Error string    `json:"error"`
	Cause string    `json:"cause"`
}

type ChoiceState struct {
	Type    StateType    `json:"type"`
	Choices []ChoiceItem `json:"choices"`
	Default string       `json:"default,omitempty"`
}
type ChoiceItem struct {
	Variable                  string `json:"variable"`
	NumericEquals             *int   `json:"NumericEqual,omitempty"`
	NumericNotEquals          *int   `json:"NumericNotEqual,omitempty"`
	NumericLessThan           *int   `json:"NumericLessThan,omitempty"`
	NumericGreaterThan        *int   `json:"NumericGreaterThan,omitempty"`
	NumericLessThanOrEqual    *int   `json:"NumericLessThanOrEqual,omitempty"`
	NumericGreaterThanOrEqual *int   `json:"NumericGreaterThanOrEqual,omitempty"`

	StringEquals             *string `json:"StringEquals,omitempty"`
	StringNotEquals          *string `json:"StringNotEquals,omitempty"`
	StringLessThan           *string `json:"StringLessThan,omitempty"`
	StringGreaterThan        *string `json:"StringGreaterThan,omitempty"`
	StringLessThanOrEqual    *string `json:"StringLessThanOrEqual,omitempty"`
	StringGreaterThanOrEqual *string `json:"StringGreaterThanOrEqual,omitempty"`
	Next                     string  `json:"next"`
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
	APIVersion string `json:"apiVersion,omitempty"`

	Name    string       `json:"name"`
	Status  VersionLabel `json:"status,omitempty"`
	StartAt string       `json:"startAt"`

	States map[string]WorkState `json:"states"`

	Comment string `json:"comment,omitempty"`
}

func (w *Workflow) MarshalJSON() ([]byte, error) {
	type Alias Workflow
	return json.Marshal(
		&struct {
			*Alias
		}{
			Alias: (*Alias)(w),
		})
}
func (w *Workflow) UnMarshalJSON() (data []byte) {
	w.APIVersion = gjson.Get(string(data), "apiVersion").String()
	w.Name = gjson.Get(string(data), "name").String()
	status := gjson.Get(string(data), "status")
	if status.Exists() {
		w.Status = VersionLabel(status.String())
	}
	w.StartAt = gjson.Get(string(data), "startAt").String()
	comment := gjson.Get(string(data), "comment")
	if comment.Exists() {
		w.Comment = comment.String()
	}
	states := gjson.Get(string(data), "states")
	if states.Exists() {
		w.States = make(map[string]WorkState)
		states.ForEach(func(key, value gjson.Result) bool {
			stateType := gjson.Get(value.String(), "type").String()
			switch stateType {
			case "Task":
				var taskState TaskState
				err := json.Unmarshal([]byte(value.String()), &taskState)
				if err != nil {
					return false
				}
				w.States[key.String()] = taskState
			case "Choice":
				var choiceState ChoiceState
				err := json.Unmarshal([]byte(value.String()), &choiceState)
				if err != nil {
					return false
				}
				w.States[key.String()] = choiceState
			case "Fail":
				var failState FailState
				err := json.Unmarshal([]byte(value.String()), &failState)
				if err != nil {
					return false
				}
			}
			return true
		})
	}

	return nil
}
