package serverless

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/workflow"
	"testing"
)

func GenerateWorkflow() core.Workflow {
	value := 5
	example := core.Workflow{
		APIVersion: "v1",
		Comment:    "An example of basic workflow.",
		StartAt:    "getsum",
		States: map[string]core.WorkState{
			"getsum": core.TaskState{
				Type:      "Task",
				InputPath: "$.x,$.y",
				Next:      "judgesum",
			},
			"judgesum": core.ChoiceState{
				Type: "Choice",
				Choices: []core.ChoiceItem{
					{
						Variable:           "$.z",
						NumericGreaterThan: &value,
						Next:               "printsum",
					},
					{
						Variable:        "$.z",
						NumericLessThan: &value,
						Next:            "getdiff",
					},
				},
				Default: "printerror",
			},
			"printsum": core.TaskState{
				Type:       "Task",
				InputPath:  "$.z",
				ResultPath: "$.str",
				End:        true,
			},
			"getdiff": core.TaskState{
				Type:      "Task",
				InputPath: "$.x,$.y,$.z",
				Next:      "printdiff",
			},
			"printdiff": core.TaskState{
				Type:       "Task",
				InputPath:  "$.z",
				ResultPath: "$.str",
				End:        true,
			},
			"printerror": core.FailState{
				Type:  "Fail",
				Error: "DefaultStateError",
				Cause: "No Matches!",
			},
		},
	}
	return example
}

func TestExecuteWorkFlow(t *testing.T) {
	//monkey.Patch(activator.TriggerFunc, func(string, []byte) ([]byte, error) {
	//	return []byte(`{"z":3, "x":4, "y":5, "str": "hello world"}`), nil
	//})
	workflowExample := GenerateWorkflow()
	params := []byte(`{"x": 4, "y": 5}`)
	result, err := workflow.ExecuteWorkflow(&workflowExample, params)
	if err != nil {
		t.Logf("ExecuteWorkFlow failed, error: %s", err)
	}
	t.Logf("result: %s", string(result))
}
