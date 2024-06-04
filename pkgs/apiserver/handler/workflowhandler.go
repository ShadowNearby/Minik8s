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

// CreateWorkflowHandler POST /api/v1/workflows
func CreateWorkflowHandler(c *gin.Context) {
	var workflow core.Workflow
	err := c.Bind(&workflow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect workflow type"})
		return
	}
	key := fmt.Sprintf("/workflows/object/%s", workflow.Name)
	err = storage.Put(key, workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error put object"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// GetWorkflowHandler GET /api/v1/workflows/:name
func GetWorkflowHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect name in get request"})
		return
	}
	key := fmt.Sprintf("/workflows/object/%s", name)
	var workflow core.Workflow
	err := storage.Get(key, &workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error get object"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(workflow)})
}

// GetWorkflowListHandler GET /api/v1/workflows
func GetWorkflowListHandler(c *gin.Context) {
	key := "/workflows/object/"
	var workflows []core.Workflow
	err := storage.RangeGet(key, &workflows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error get object"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(workflows)})
}

// DeleteWorkflowHandler DELETE /api/v1/workflows/:name
func DeleteWorkflowHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect workflow name"})
		return
	}
	key := fmt.Sprintf("/workflows/object/%s", name)
	err := storage.Del(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error delete object"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// UpdateWorkflowHandler PUT /api/v1/workflows/:name
func UpdateWorkflowHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect workflow name"})
		return
	}
	key := fmt.Sprintf("/workflow/object/%s", name)
	var oldWorkflow core.Workflow
	err := storage.Get(key, &oldWorkflow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "donot have workflow before"})
		return
	}
	var newWorkflow core.Workflow
	err = c.Bind(&newWorkflow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect workflow tyep"})
		return
	}
	storage.Put(key, newWorkflow)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// TriggerWorkflowHandler POST /api/v1/workflows/:name/trigger
func TriggerWorkflowHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect name"})
		return
	}
	var request core.WorkFlowTriggerRequest
	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect workflow reqeust type"})
		return
	}
	logger.Infof("request name: %s, params: %s", request.Name, request.Params)
	storage.RedisInstance.PublishMessage(constants.ChannelWorkflowTrigger, utils.JsonMarshal(request))
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
