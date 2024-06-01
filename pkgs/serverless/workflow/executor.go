package workflow

import (
	"errors"
	"fmt"
	"math"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/serverless/activator"
	"minik8s/utils"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const epsilon = 10e-4

func CheckNode(nodeName string) bool {
	url := fmt.Sprintf("http://%s:%s/api/v1/functions/%s", config.ClusterMasterIP, config.ApiServerPort, nodeName)
	_, _, err := utils.SendRequest("GET", url, nil)
	log.Error("CheckNode ", err)
	if err != nil {
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

func ExecuteWorkflow(workflow *core.Workflow, params []byte) ([]byte, error) {
	startNode := workflow.StartAt
	if startNode == "" {
		log.Error("[ExecuteWorkflow] the workflow start node is not exist")
		return nil, errors.New("workflow start node is empty")
	}
	currentNode, ok := workflow.States[startNode]
	if !ok {
		log.Error("[ExecuteWorkflow] the workflow states is not exist")
		return nil, errors.New("workflow states is empty")
	}
	currentNodeName := startNode
	log.Info("[ExecuteWorkflow] current node name is: ", currentNodeName, " with params: ", string(params), " and type ", reflect.TypeOf(currentNode))

	for {
		prevNodeName := currentNode
		err := error(nil)

		switch reflect.TypeOf(currentNode).Name() {
		case string(core.Task):
			{
				params, err = ExecuteTask(currentNode.(core.TaskState), currentNodeName, params)
				if err != nil {
					log.Error("[ExecuteWorkflow] task execution failed, the current node is ", currentNodeName)
					return nil, err
				}

			}
		case string(core.Choice):
			{
				currentNodeName, err = ExecuteChoice(currentNode.(core.ChoiceState), params)
				if err != nil {
					log.Error("[ExecuteWorkflow] choice execution failed, the current node is ", currentNodeName)
					return nil, err
				}
			}
		case string(core.Fail):
			{
				result := ExecuteFail(currentNode.(core.FailState))
				return []byte(result), nil
			}
		default:
			{
				log.Info(reflect.TypeOf(currentNode).Name())
				return nil, errors.New("the current node's type is invalid")
			}
		}
		currentNode = workflow.States[currentNodeName]

		// don't allow loop now
		if currentNode == nil || prevNodeName == currentNodeName {
			break
		}
	}
	return []byte(currentNodeName), errors.New("the workflow is not valid")
}
func replaceSingleQuotesWithDoubleQuotes(str string) string {
	// the default string in dict is single quotes, need to replace it with double quotes
	bytes := []byte(str)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == '\'' {
			bytes[i] = '"'
		}
	}
	return string(bytes)
}

func ExecuteTask(task core.TaskState, functionName string, params []byte) ([]byte, error) {

	if functionName == "" {
		return nil, errors.New("task resource is empty")
	}

	// check the function is valid or not
	if !CheckNode(functionName) {
		return nil, errors.New("function is not valid")
	}
	// try to trigger the function
	// if the InputPath is not empty, need to parse the params to abstract the input
	inputParams := params
	err := error(nil)
	if task.InputPath != "" {
		inputParams, err = ParseParams(params, task.InputPath)
		if err != nil {
			return nil, err
		}
	}
	log.Info("Triggering Function ======>")
	result, err := activator.TriggerFunc(functionName, inputParams)
	log.Info("======> Success")
	if err != nil {
		return nil, err
	}
	log.Info(result)

	// python's dict is single quotes, need to replace it with double quotes

	paramsStr := string(result)
	paramsStr = replaceSingleQuotesWithDoubleQuotes(paramsStr)
	result = []byte(paramsStr)
	log.Info(paramsStr)

	// if the ResultPath is not empty, need to parse the result to abstract the output
	if task.ResultPath != "" {
		result, err = ParseParams(result, task.ResultPath)

		if err != nil {
			return nil, err
		}
	}

	return result, nil
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

// isNumeric check whether the variable's type is numeric
func isNumeric(variable interface{}) bool {
	switch variable.(type) {
	// actually, if use gjson to get the value, the type is float64 default
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
		return true
	default:
		return false
	}
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
