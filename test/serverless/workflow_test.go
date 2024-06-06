package serverless

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/workflow"
	"minik8s/utils"
	"os"
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

func TestCreateWorkflow(t *testing.T) {
	url := fmt.Sprintf("http://%s:8090/api/v1/workflows/", config.ClusterMasterIP)
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "function/workflow.json"))
	if err != nil {
		t.Errorf("read file error")
		return
	}
	var workflow core.Workflow
	_ = utils.JsonUnMarshal(string(file), &workflow)
	code, info, err := utils.SendRequest("POST", url, file)
	if err != nil || code != 200 {
		var infotype core.InfoType
		utils.JsonUnMarshal(info, &infotype)
		t.Errorf("send create workflow error: %s", infotype.Error)
		return
	}
}

func TestTriggerWorkflow(t *testing.T) {
	url := fmt.Sprintf("http://%s:8090/api/v1/workflows/%s/trigger", config.ClusterMasterIP, "app-yolo")
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "serverless/common/trigger.json"))
	if err != nil {
		t.Errorf("read file error")
		return
	}
	var trigger core.TriggerMessage
	_ = utils.JsonUnMarshal(string(file), &trigger)
	code, info, err := utils.SendRequest("POST", url, file)
	if err != nil || code != 200 {
		var infotype core.InfoType
		utils.JsonUnMarshal(info, &infotype)
		t.Errorf("send create workflow error: %s", infotype.Error)
		return
	}
}

func TestWorkflowParse(t *testing.T) {
	states := make(map[string]core.WorkState)

	taskState := core.TaskState{
		Type:       core.Task,
		InputPath:  "input/path",
		ResultPath: "result/path",
		Next:       "NextState",
		End:        false,
	}
	states["task"] = taskState

	failState := core.FailState{
		Type:  core.Fail,
		Error: "Some error",
		Cause: "Some cause",
	}
	states["fail"] = failState

	choiceState := core.ChoiceState{
		Type:    core.Choice,
		Choices: []core.ChoiceItem{
			// Your ChoiceItem initialization here
		},
		Default: "DefaultChoice",
	}
	states["choice"] = choiceState

	// Retrieve and identify state
	for key, state := range states {
		switch s := state.(type) {
		case core.TaskState:
			fmt.Printf("Key: %s, Type: TaskState, Value: %+v\n", key, s)
		case core.FailState:
			fmt.Printf("Key: %s, Type: FailState, Value: %+v\n", key, s)
		case core.ChoiceState:
			fmt.Printf("Key: %s, Type: ChoiceState, Value: %+v\n", key, s)
		default:
			fmt.Printf("Key: %s, Type: Unknown\n", key)
		}
	}
}
