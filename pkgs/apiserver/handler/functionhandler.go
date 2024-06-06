package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func FunctionKeyPrefix(name string) string {
	return fmt.Sprintf("/functions/object/%s", name)
}

// CreateFunctionHandler POST /api/v1/functions
func CreateFunctionHandler(c *gin.Context) {
	var function core.Function
	err := c.Bind(&function)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check the parameters
	if function.Name == "" {
		logger.Errorf("Function name is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "function is empty"})
		return
	}
	if function.Path == "" {
		logger.Errorf("Funtion path is empty")
		c.JSON(http.StatusBadRequest, gin.H{"create": "function path is empty"})
		return
	}
	key := FunctionKeyPrefix(function.Name)
	err = storage.Put(key, function)
	if err != nil {
		logger.Error("[FunctionCreateHandler] error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelCreate), utils.JsonMarshal(function))
	c.JSON(http.StatusOK, gin.H{"data": "create function success"})
}

// GetFunctionHandler GET /api/v1/functions/:name
func GetFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "function name is empty"})
		return
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
	err = storage.Del(functionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete function"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelDelete), function.Name)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// UpdateFunctionHandler PUT /api/v1/functions/:name
func UpdateFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	var oldFunc, newFunc core.Function
	key := fmt.Sprintf("/functions/object/%s", name)
	logger.Infof(key)
	err := storage.Get(key, &oldFunc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "create before update"})
		return
	}
	err = c.Bind(&newFunc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect function type"})
		return
	}
	err = storage.Put(key, newFunc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot set object"})
		return
	}
	functions := []core.Function{oldFunc, newFunc}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelFunction, constants.ChannelUpdate), utils.JsonMarshal(functions))
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// TriggerFunctionHandler POST /api/v1/functions/:name/trigger
// create an ID to identify the result, client will use this id to ask for result
func TriggerFunctionHandler(c *gin.Context) {
	name := c.Param("name")
	functionName := FunctionKeyPrefix(name)
	var function core.Function
	if err := storage.Get(functionName, &function); err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}

	triggerMsg := core.TriggerMessage{}
	err := c.Bind(&triggerMsg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect trigger message type"})
		return
	}
	triggerMsg.ID = utils.GenerateUUID(4)
	storage.RedisInstance.PublishMessage(constants.ChannelFunctionTrigger, utils.JsonMarshal(triggerMsg))
	c.JSON(http.StatusOK, gin.H{"data": triggerMsg.ID})
}

// GetAllFunctionsHandler GET /api/v1/functions
func GetAllFunctionsHandler(c *gin.Context) {
	var functionConfigs []core.Function
	err := storage.RangeGet("/functions/object", &functionConfigs)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(functionConfigs)})
}

// CreateTaskHandler /api/v1/functions/task POST
func CreateTaskHandler(c *gin.Context) {
	var task core.PingSource
	err := c.Bind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect Ping source type"})
		return
	}
	if task.MetaData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task name cannot be empty"})
		return
	}
	if task.Spec.Sink.Ref.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trigger function name cannot be empty"})
		return
	}
	err = storage.Put(fmt.Sprintf("/functions/tasks/%s", task.MetaData.Name), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save task"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelTask, constants.ChannelCreate), utils.JsonMarshal(task))
	c.JSON(http.StatusOK, gin.H{"data": "ok"})

}

// UpdateTaskHandler /api/v1/functions/task/:name "POST"
func UpdateTaskHandler(c *gin.Context) {
	var task core.PingSource
	var newTask core.PingSource
	name := c.Param("name")
	err := c.Bind(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect ping source type"})
		return
	}
	if newTask.MetaData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task name cannot be empty"})
		return
	}
	if newTask.Spec.Sink.Ref.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trigger function name cannot be empty"})
		return
	}
	err = storage.Get(fmt.Sprintf("/functions/tasks/%s", name), &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update task does not exist"})
		return
	}
	err = storage.Put(fmt.Sprintf("/functions/tasks/%s", task.MetaData.Name), newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save task"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelTask, constants.ChannelUpdate), utils.JsonMarshal(newTask))
}

// DeleteTaskHandler /api/v1/functions/task/name "DELETE"
func DeleteTaskHandler(c *gin.Context) {
	name := c.Param("name")
	var oldTask core.PingSource
	err := storage.Get(fmt.Sprintf("/functions/tasks/%s", name), &oldTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete task does not exist"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelTask, constants.ChannelDelete), name)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// GetTaskHandler /api/v1/functions/task/:name "GET"
func GetTaskHandler(c *gin.Context) {
	name := c.Param("name")
	var task core.PingSource
	err := storage.Get(fmt.Sprintf("/functions/tasks/%s", name), &task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot get task does not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(task)})
}

// GetAllTaskHandler /api/v1/functions/task "GET"
// this function should be called every time apiserver restart to save data in taskHandler
func GetAllTaskHandler(c *gin.Context) {
	var tasks []core.PingSource
	err := storage.RangeGet("/functions/tasks", &tasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(tasks)})
}

// GetTriggerResult /api/v1/functions/result/:id "GET"
// will delete result after get once
func GetTriggerResult(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trigger id should not be null"})
		return
	}
	key := fmt.Sprintf("/functions/result/%s", id)
	var result core.TriggerResult
	err := storage.Get(key, &result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "function trigger has not done yet"})
		return
	}
	// to reduce redundant keys, we delete data
	storage.Del(key)
	c.JSON(http.StatusOK, gin.H{"data": result.Result})
}

// SetTriggerResult /api/v1/functions/result "POST"
func SetTriggerResult(c *gin.Context) {
	var result core.TriggerResult
	err := c.Bind(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect result type"})
		return
	}
	key := fmt.Sprintf("/functions/result/%s", result.ID)
	err = storage.Put(key, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error put result"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
