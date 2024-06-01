package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"
)

func FunctionKeyPrefix(name string) string {
	return fmt.Sprintf("/registry/functions/%s", name)
}

// CreateFunctionHandler POST /api/v1/functions
func CreateFunctionHandler(c *gin.Context) {
	var function core.Function
	err := c.Bind(&function)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Info("[FunctionCreateHandler] function: ", function)
	// check the parameters
	if function.Name == "" {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "function is empty"})
		return
	}
	if function.Path == "" {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"create": "function path is empty"})
		return
	}
	key := FunctionKeyPrefix(function.Name)
	err = storage.Put(key, function)
	logger.Info("[FunctionCreateHandler] init function finished")
	if err != nil {
		logger.Error("[FunctionCreateHandler] error: ", err.Error())
		c.JSON(http.StatusInternalServerError, []byte("create: "+err.Error()))
	} else {
		logger.Info("[FunctionCreateHandler] success")
		c.JSON(http.StatusOK, []byte("create: "+"function create success"))
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelCreate), function)
	newFunction := []core.Function{core.Function{}, function}
	storage.RedisInstance.PublishMessage(constants.ChannelPodSchedule, utils.JsonMarshal(newFunction))
	if err != nil {
		logger.Error("[FunctionCreateHandler] error: ", err.Error())
		c.JSON(http.StatusInternalServerError, []byte("create: "+err.Error()))
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// GetFunctionHandler GET /api/v1/functions/:name
func GetFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "function name is empty"})
	}
	var function core.Function
	functionName := FunctionKeyPrefix(name)
	err := storage.Get(functionName, &function)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get function"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(function)})
}

// DeleteFunctionHandler DELETE /api/v1/functions/:name
func DeleteFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "function name is empty"})
	}
	var function core.Function
	functionName := FunctionKeyPrefix(name)
	err := storage.Get(functionName, &function)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	function.Status = core.DELETE
	err = storage.Del(functionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete function"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(function)})
}

// UpdateFunctionHandler PUT /api/v1/functions/:name
func UpdateFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "function name is empty"})
	}
	functionConfig := &core.Function{}
	if err := c.Bind(functionConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	functionName := FunctionKeyPrefix(name)
	preFunctionConfig := &core.Function{}
	if err := storage.Get(functionName, &preFunctionConfig); err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	if err := storage.Put(functionName, &preFunctionConfig); err != nil {
		logger.Errorf("function %s not found", name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "function put error"})
		return
	}
	////TODO : adjust update
	functions := []core.Function{*preFunctionConfig, *functionConfig}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelUpdate), utils.JsonMarshal(functions))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// TriggerFunctionHandler POST /api/v1/functions/:name/trigger
func TriggerFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	functionName := FunctionKeyPrefix(name)
	var function core.Function
	if err := storage.Get(functionName, &function); err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}

	paramsRaw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request, err := json.Marshal(
		struct {
			Name   string          `json:"name"`
			Params json.RawMessage `json:"params"`
		}{
			Name:   name,
			Params: paramsRaw,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Info("[TriggerFunctionHandler] function: ", functionName)
	err = storage.Put(functionName, string(request))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	logger.Info("[TriggerFunctionHandler] function: ", functionName)
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelTrigger), utils.JsonMarshal(request))

	if err != nil {
		// send request error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

}

// GetAllFunctionsHandler GET /api/v1/functions
func GetAllFunctionsHandler(c *gin.Context) {
	var functionConfigs []core.Pod
	err := storage.RangeGet("/registry/functions/", &functionConfigs)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(functionConfigs)})

}
