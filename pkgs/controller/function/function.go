package function

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/serverless/activator"
	"minik8s/utils"

	logger "github.com/sirupsen/logrus"
)

type FuncController struct {
}

func (f *FuncController) GetChannel() string {
	return constants.ChannelFunction
}

func (f *FuncController) ListenOtherChannels() {
	httpTriggerMessage := storage.RedisInstance.SubscribeChannel(constants.ChannelFunctionTrigger)
	for message := range httpTriggerMessage {
		err := f.HandleHttpTrigger(message.Payload)
		if err != nil {
			logger.Errorf("handle trigger error: %s", err.Error())
		}
	}
}
func (f *FuncController) HandleCreate(message string) error {
	var fnc core.Function
	err := utils.JsonUnMarshal(message, &fnc)
	if err != nil {
		logger.Errorf("unmarshal function error: %s", err.Error())
		return err
	}
	err = activator.InitFunction(fnc.Name, fnc.Path)
	if err != nil {
		logger.Errorf("init function error: %s", err.Error())
	}
	return err
}

func (f *FuncController) HandleUpdate(message string) error {
	// we don't support update
	functions := []core.Function{}
	err := utils.JsonUnMarshal(message, &functions)
	if err != nil {
		logger.Error("unmarshal functions error")
		return err
	}
	err = f.HandleDelete(functions[0].Name)
	if err != nil {
		logger.Error("delete functions error")
		return err
	}
	err = f.HandleCreate(utils.JsonMarshal(functions[1]))
	if err != nil {
		logger.Error("create functions error")
		return err
	}
	return nil
}

func (f *FuncController) HandleDelete(message string) error {
	// message is function name
	err := activator.DeleteFunc(message)
	if err != nil {
		logger.Errorf("delete function error: %s", err.Error())
		return err
	}
	return nil
}

func (f *FuncController) HandleHttpTrigger(message string) error {
	// message is TriggerMessage
	var triggerMessage core.TriggerMessage
	utils.JsonUnMarshal(message, &triggerMessage)
	result, err := activator.TriggerFunc(triggerMessage.Name, []byte(triggerMessage.Params))
	if err != nil {
		logger.Errorf("trigger function error: %s", err.Error())
	}
	logger.Infof("trigger result: %s", string(result))
	return err
}
