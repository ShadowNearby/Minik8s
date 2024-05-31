package function

import (
	"encoding/json"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
)

type FunctionController struct{}

func (f *FunctionController) GetChannel() string {
	return constants.ChannelFunction
}
func (f *FunctionController) HandleCreate(message string) error {
	var function core.Function
	err := json.Unmarshal([]byte(message), &function)
	if err != nil {
		return err
	}
	return nil
}
func (f *FunctionController) HandleUpdate(message string) error {
	return nil
}
func (f *FunctionController) HandleDelete(message string) error {
	return nil
}
func (f *FunctionController) HandleTrigger(message string) error {
	return nil
}
