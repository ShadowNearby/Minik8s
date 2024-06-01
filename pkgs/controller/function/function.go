package function

import (
	"encoding/json"
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/serverless/activator"
)

type FunctionController struct{}

func (f *FunctionController) GetChannel() string {
	return constants.ChannelFunction
}
func (f *FunctionController) HandleCreate(message string) error {
	var function core.Function
	err := json.Unmarshal([]byte(message), &function)
	logger.Info("[FunctionCreateHandler] function: ", function)
	// check the parameters
	if function.Name == "" {
		logger.Errorf("create: " + "function name is empty")
		return errors.New("create: " + "function name is empty")
	}
	if function.Path == "" {
		logger.Errorf("create: " + "function path is empty")
		return errors.New("create: " + "function path is empty")
	}
	err = activator.InitFunc(function.Name, function.Path)
	logger.Info("[FunctionCreateHandler] init function finished")
	if err != nil {
		logger.Error("[FunctionCreateHandler] error: ", err.Error())
		return err
	} else {
		logger.Info("[FunctionCreateHandler] success")
		return nil
	}
}
func (f *FunctionController) HandleUpdate(message string) error {
	function := &core.Function{}
	err := json.Unmarshal([]byte(message), function)
	if err != nil {
		logger.Error("[FunctionUpdateHandler] function: ", function)
	}
	logger.Info("[FunctionUpdateHandler] function: ", function)

	// delete the old function and create the new function
	err = activator.DeleteFunc(function.Name)
	if err != nil {
		logger.Errorf("update: " + err.Error())
		return errors.New("update: " + err.Error())
	}
	logger.Info("[FunctionUpdateHandler] delete function finished")
	err = activator.InitFunc(function.Name, function.Path)
	logger.Info("[FunctionUpdateHandler] update function finished")
	if err != nil {
		logger.Errorf("update: " + err.Error())
		return errors.New("update: " + err.Error())
	} else {
		logger.Errorf("update: function update success")
		return nil
	}
}
func (f *FunctionController) HandleDelete(message string) error {
	function := &core.Function{}
	err := json.Unmarshal([]byte(message), function)
	logger.Info("[FunctionDeleteHandler] function: ", function)
	if err != nil {
		logger.Errorf("delete: " + err.Error())
		return errors.New("delete: " + err.Error())
	}
	// check the parameters
	if function.Name == "" {
		logger.Errorf("delete: " + "function name is empty")
		return errors.New("function name is empty")
	}
	err = activator.DeleteFunc(function.Name)
	logger.Info("[FunctionDeleteHandler] delete function finished")
	if err != nil {
		logger.Errorf("delete: " + err.Error())
		return errors.New("delete: " + err.Error())
	} else {
		logger.Errorf("delete: " + "function delete success")
		return nil
	}
}
func (f *FunctionController) HandleTrigger(message string) error {
	nameField := gjson.Get(message, "name")
	if !nameField.Exists() {
		logger.Errorf("execute: " + "function name is empty")
		return errors.New("function name is empty")
	}
	name := nameField.String()
	paramsField := gjson.Get(message, "params")
	if !paramsField.Exists() {
		logger.Errorf("execute: " + "function params is empty")
		return errors.New("function params is empty")
	}
	params := paramsField.String()
	logger.Info("[FunctionTriggerHandler] name: ", name, ", params: ", params)
	result, err := activator.TriggerFunc(name, []byte(params))
	if err != nil {
		logger.Errorf("execute: " + err.Error())
		return errors.New("execute: " + err.Error())
	}
	logger.Info("execute: " + string(result))
	return nil
}
