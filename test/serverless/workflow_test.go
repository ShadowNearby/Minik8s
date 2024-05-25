package serverless

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/workflow"
	"testing"
)

//func TestGetParam(t *testing.T) {
//	data := `"{"z": 9, "x": 5, "y": 4}"`
//	path := "$.x"
//	result := gjson.Get(data, path[2:])
//	if !result.Exists() {
//		t.Errorf("GetParam failed, error: %s", "result is not exist")
//	}
//}
//func TestWorkflow(t *testing.T) {
//	// test1: only integer, nothing need to ignore
//	params := []byte(`{"x": 4, "y": 5}`)
//	inputPath := "$.x,$.y"
//	result, err := workflow.ParseParams(params, inputPath)
//	if err != nil {
//		t.Errorf("ParseParams failed, error: %s", err)
//	}
//	t.Logf("result: %s", string(result))
//
//	// test2: only integer, but need to ignore
//	params = []byte(`{"x": 4, "y": 5}`)
//	inputPath = "$.x"
//	result, err = workflow.ParseParams(params, inputPath)
//	if err != nil {
//		t.Errorf("ParseParams failed, error: %s", err)
//	}
//	t.Logf("result: %s", string(result))
//
//	// test3: integer and string, nothing need to ignore
//	params = []byte(`{"x": 4, "str": "hello"}`)
//	inputPath = "$.x,$.str"
//	result, err = workflow.ParseParams(params, inputPath)
//	if err != nil {
//		t.Errorf("ParseParams failed, error: %s", err)
//	}
//	t.Logf("result: %s", string(result))
//
//	// test4: integer and string, but need to ignore
//	params = []byte(`{"x": 4, "str": "hello"}`)
//	inputPath = "$.x"
//	result, err = workflow.ParseParams(params, inputPath)
//	if err != nil {
//		t.Errorf("ParseParams failed, error: %s", err)
//	}
//	t.Logf("result: %s", string(result))
//}
//
//func TestHasField(t *testing.T) {
//	value := 5
//	chElem := core.ChoiceItem{
//		Variable:           "$.z",
//		NumericGreaterThan: &value,
//		Next:               "PrintSum",
//	}
//	if !workflow.HasField(chElem, "NumericGreaterThan") {
//		t.Errorf("HasField failed, error: %s, %d", "NumericGreaterThan", chElem.NumericGreaterThan)
//	}
//	if workflow.HasField(chElem, "NumericLessThan") {
//		t.Errorf("HasField failed, error: %s, %d", "NumericLessThan", chElem.NumericLessThan)
//	}
//}
//
//func TestExecuteChoice(t *testing.T) {
//	value1 := 1
//	value2 := 2
//	choice := core.ChoiceState{
//		Type: "Choice",
//		Choices: []core.ChoiceItem{
//			{
//				Variable:      "$.foo",
//				NumericEquals: &value1,
//				Next:          "FirstMatchState",
//			},
//			{
//				Variable:      "$.foo",
//				NumericEquals: &value2,
//				Next:          "SecondMatchState",
//			},
//		},
//		Default: "DefaultState",
//	}
//
//	params := []byte(`{"foo": 1}`)
//	next, err := workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	t.Logf("next: %s", next)
//
//	params = []byte(`{"foo": 2}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//
//	// the second test
//	value := 5
//	choice = core.ChoiceState{
//		Type: "Choice",
//		Choices: []core.ChoiceItem{
//			{
//				Variable:           "$.z",
//				NumericGreaterThan: &value,
//				Next:               "PrintSum",
//			},
//			{
//				Variable:        "$.z",
//				NumericLessThan: &value,
//				Next:            "GetDiff",
//			},
//		},
//		Default: "PrintError",
//	}
//
//	params = []byte(`{"z": 6}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "PrintSum" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//
//	params = []byte(`{"z": 4}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "GetDiff" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//
//	params = []byte(`{"z": 5}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "PrintError" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//
//	// the third test, check string
//	strValue := "hello"
//	choice = core.ChoiceState{
//		Type: "Choice",
//		Choices: []core.ChoiceItem{
//			{
//				Variable:     "$.str",
//				StringEquals: &strValue,
//				Next:         "PrintSum",
//			},
//			{
//				Variable:          "$.str",
//				StringGreaterThan: &strValue,
//				Next:              "PrintSum",
//			},
//			{
//				Variable:       "$.str",
//				StringLessThan: &strValue,
//				Next:           "GetDiff",
//			},
//		},
//		Default: "PrintError",
//	}
//
//	params = []byte(`{"str": "hello"}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "PrintSum" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//
//	params = []byte(`{"str": "world"}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "PrintSum" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//
//	params = []byte(`{"str": "hell"}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "GetDiff" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//
//	params = []byte(`{"str": "world hello"}`)
//	next, err = workflow.ExecuteChoice(choice, params)
//	if err != nil {
//		t.Errorf("ExecuteChoice failed, error: %s", err)
//	}
//	if next != "PrintSum" {
//		t.Errorf("ExecuteChoice failed, error: %s", next)
//	}
//}

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
