package eventfilter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/serverless/activator"
	"time"
)

func FunctionSync(target string) {
	// establish websocket connection
	for {
		err := connect(target)
		if err != nil {
			log.Error("[FunctionSync] WebSocket connect fail: ", err)
		}
		time.Sleep(5 * time.Second) // wait 5 seconds to reconnect
	}
}

func connect(target string) error {
	// establish websocket connection
	createChannel := constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelCreate)
	updateChannel := constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelUpdate)
	deleteChannel := constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelDelete)
	functionChannel := storage.RedisInstance.SubscribeChannel(createChannel)
	functionChannel = storage.RedisInstance.SubscribeChannel(updateChannel)
	functionChannel = storage.RedisInstance.SubscribeChannel(deleteChannel)
	go func() {
		for message := range functionChannel {

			msg := message.Payload
			if len(msg) == 0 {
				continue
			}
			fmt.Printf("[client %s] %s\n", target, message)
			op := gjson.Get(msg, "status")
			// function trigger
			if !op.Exists() {
				go FunctionTriggerHandler(msg)
				continue
			}
			switch op.String() {
			case "create":
				{
					go FunctionCreateHandler(msg)
				}
			case "delete":
				{
					go FunctionDeleteHandler(msg)
				}
			case "update":
				{
					go FunctionUpdateHandler(msg)
				}
			}
		}
	}()
	return nil
}

// FunctionTriggerHandler the trigger format: {"name": "function name", "params": "function params"}
func FunctionTriggerHandler(message string) {
	nameField := gjson.Get(message, "name")
	if !nameField.Exists() {
		log.Errorf("execute: " + "function name is empty")
		return
	}

	name := nameField.String()
	paramsField := gjson.Get(message, "params")
	if !paramsField.Exists() {
		log.Errorf("execute: " + "function params is empty")
		return
	}

	params := paramsField.String()
	log.Info("[FunctionTriggerHandler] name: ", name, ", params: ", params)
	result, err := activator.TriggerFunc(name, []byte(params))
	if err != nil {
		log.Errorf("execute: " + err.Error())
		return
	}

	log.Info("execute: " + string(result))
}

func FunctionDeleteHandler(message string) {
	function := &core.Function{}
	function.UnMarshalJSON([]byte(message))
	log.Info("[FunctionDeleteHandler] function: ", function)

	// check the parameters
	if function.Name == "" {
		log.Errorf("delete: " + "function name is empty")
		return
	}

	err := activator.DeleteFunc(function.Name)
	log.Info("[FunctionDeleteHandler] delete function finished")
	if err != nil {
		log.Errorf("delete: " + err.Error())
	} else {
		log.Errorf("delete: " + "function delete success")
	}

}

func FunctionUpdateHandler(message string) {
	function := &core.Function{}
	function.UnMarshalJSON([]byte(message))
	log.Info("[FunctionUpdateHandler] function: ", function)

	// delete the old function and create the new function
	err := activator.DeleteFunc(function.Name)
	if err != nil {
		log.Errorf("update: " + err.Error())
		return
	}
	log.Info("[FunctionUpdateHandler] delete function finished")
	err = activator.InitFunc(function.Name, function.Path)
	log.Info("[FunctionUpdateHandler] update function finished")
	if err != nil {
		log.Errorf("update: " + err.Error())
	} else {
		log.Errorf("update: function update success")
	}

}

func FunctionCreateHandler(message string) {
	function := &core.Function{}
	function.UnMarshalJSON([]byte(message))
	log.Info("[FunctionCreateHandler] function: ", function)

	// check the parameters
	if function.Name == "" {
		log.Errorf("create: " + "function name is empty")
		return
	}
	if function.Path == "" {
		log.Errorf("create: " + "function path is empty")
	}

	err := activator.InitFunc(function.Name, function.Path)
	log.Info("[FunctionCreateHandler] init function finished")
	if err != nil {
		log.Error("[FunctionCreateHandler] error: ", err.Error())
	} else {
		log.Info("[FunctionCreateHandler] success")
	}
}
