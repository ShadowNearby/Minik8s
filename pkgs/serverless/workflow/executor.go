package workflow

import (
	"errors"
	"fmt"
	"math"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/activator"
	"minik8s/utils"
	"net/http"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const epsilon = 10e-4

func checkValidFunc(nodeName string) bool {
	url := fmt.Sprintf("http://%s:%s/api/v1/functions/%s", config.ClusterMasterIP, config.ApiServerPort, nodeName)

	code, _, err := utils.SendRequest("GET", url, nil)
	if err != nil || code != http.StatusOK {
		return false
	}
	return true
}
func ParseParams(params []byte, inputPath string) ([]byte, error) {
	wanted := strings.Split(inputPath, ",")

	filterdParams := make(map[string]interface{})
	for _, elem := range wanted {
		name := elem[2:]
		value := gjson.Get(string(params), name)
		if !value.Exists() {
			log.Error("[ParseParams] the params is not valid, the name is: ", name, ", the params is: ", string(params))
			return nil, errors.New("the params is not valid")
		}
		filterdParams[name] = value.Value()
	}

	jsonData := utils.JsonMarshal(filterdParams)
	if jsonData == "" {
		log.Error("[ParseParams] the params is not valid, the params is: ", string(params))
	}
	return []byte(jsonData), nil
}

func ExecuteWorkflow(workflow *core.Workflow, params []byte) (string, error) {
	startNode := workflow.StartAt
	if startNode == "" {
		log.Error("[ExecuteWorkflow] the workflow start node is not exist")
		return "", errors.New("workflow start node is empty")
	}
	currentState, ok := workflow.States[startNode]
	if !ok {
		log.Error("[ExecuteWorkflow] the workflow states is not exist")
		return "", errors.New("workflow states is empty")
	}
	currentNodeName := startNode
	log.Info("[ExecuteWorkflow] current node name is: ", currentNodeName, " with params: ", string(params), " and type ", reflect.TypeOf(currentState))

	for {
		prevNodeName := currentNodeName
		err := error(nil)

		switch currentState := currentState.(type) {
		case core.TaskState:
			{
				params, err = ExecuteTask(currentState, currentNodeName, params)
				if err != nil {
					log.Error("[ExecuteWorkflow] task execution failed, the current node is ", currentNodeName)
					return "", err
				}
				if currentState.End {
					return string(params), nil
				}
				currentNodeName = currentState.Next

			}
		case core.ChoiceState:
			{
				currentNodeName, err = ExecuteChoice(currentState, params)
				if err != nil {
					log.Error("[ExecuteWorkflow] choice execution failed, the current node is ", currentNodeName)
					return "", err
				}
			}
		case core.FailState:
			{
				result := ExecuteFail(currentState)
				return result, nil
			}
		default:
			{
				log.Info(reflect.TypeOf(currentState).Name())
				return "", errors.New("the current node's type is invalid")
			}
		}
		currentState = workflow.States[currentNodeName]

		// don't allow loop now
		if currentState == nil || prevNodeName == currentNodeName {
			log.Error("current state is nil or in loop")
			break
		}
	}
	return currentNodeName, errors.New("the workflow is not valid")
}

func ExecuteTask(task core.TaskState, functionName string, params []byte) ([]byte, error) {
	log.Info("execute task")
	if functionName == "" {
		return nil, errors.New("task resource is empty")
	}

	// check the function is valid or not
	if !checkValidFunc(functionName) {
		return nil, errors.New("function is not valid")
	}
	inputParams := params
	err := error(nil)
	if task.InputPath != "" {
		inputParams, err = ParseParams(params, task.InputPath)
		if err != nil {
			return nil, err
		}
	}
	result, err := activator.TriggerFunc(functionName, inputParams)
	return []byte(result), err
}
func ExecuteFail(fail core.FailState) string {
	result := fmt.Sprintf("Fail: %s, Cause: %s", fail.Error, fail.Cause)
	log.Error("Fail: ", fail.Error, "Case: ", fail.Cause)
	return result
}
func HasField(obj interface{}, fieldName string) bool {
	t := reflect.ValueOf(obj)
	value := t.FieldByName(fieldName)
	return value.Kind() == reflect.Ptr && !value.IsNil()
}

// isString check whether the variable's type is string
func isString(variable interface{}) bool {
	switch variable.(type) {
	case string:
		return true
	default:
		return false
	}
}

func ExecuteChoice(choiceState core.ChoiceState, params []byte) (string, error) {
	log.Info("execute choice")
	for _, choice := range choiceState.Choices {
		variable := gjson.Get(string(params), choice.Variable[2:])
		if !variable.Exists() {
			return "", errors.New("the variable is not exist")
		}
		value := variable.Value()
		if HasField(choice, "NumericEquals") {
			val, ok := value.(float64)
			if ok && math.Abs(float64(*choice.NumericEquals)-val) < epsilon {
				log.Info("[ExecuteChoice] NumericEquals: ", val, " = ", *choice.NumericEquals, " return")
				return choice.Next, nil
			}
		} else if HasField(choice, "StringEquals") {
			if isString(value) && *choice.StringEquals == value {
				return choice.Next, nil
			}
		} else if HasField(choice, "NumericNotEquals") {
			val, ok := value.(float64)
			if ok && math.Abs(float64(*choice.NumericNotEquals)-val) > epsilon {
				return choice.Next, nil
			}
		} else if HasField(choice, "StringNotEquals") {
			if isString(value) && *choice.StringNotEquals != value {
				return choice.Next, nil
			}
		} else if HasField(choice, "NumericLessThan") {
			val, ok := value.(float64)
			if ok && float64(*choice.NumericLessThan) > val {
				return choice.Next, nil
			}
		} else if HasField(choice, "StringLessThan") {
			val, ok := value.(string)
			if ok && *choice.StringLessThan > val {
				return choice.Next, nil
			}
		} else if HasField(choice, "NumericGreaterThan") {
			val, ok := value.(float64)
			if ok && float64(*choice.NumericGreaterThan) < val {
				return choice.Next, nil
			}
		} else if HasField(choice, "StringGreaterThan") {
			val, ok := value.(string)
			if ok && *choice.StringGreaterThan < val {
				return choice.Next, nil
			}
		}

	}
	log.Info("[ExecuteChoice] default: ", choiceState.Default)
	if choiceState.Default != "" {
		return choiceState.Default, nil
	}
	return "", errors.New("the choice is not valid")
}
