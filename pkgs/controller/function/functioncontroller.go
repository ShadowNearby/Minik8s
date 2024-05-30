package function

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
)

type FunctionController struct{}

func (f *FunctionController) GetChannel() string {

	return constants.ChannelFunction
}
func (f *FunctionController) HandleCreate(info string) error {
	var function core.Function
	err := json.Unmarshal([]byte(info), &function)
	if err != nil {
		return err
	}
	return nil
}
func (f *FunctionController) HandleGet(c *gin.Context) {

}
func (f *FunctionController) HandleUpdate(c *gin.Context) {

}
func (f *FunctionController) HandleDelete(c *gin.Context) {

}
