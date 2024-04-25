package handler

import "github.com/gin-gonic/gin"

// CreateWorkflowHandler POST /api/v1/workflows
func CreateWorkflowHandler(c *gin.Context) {}

// GetWorkflowHandler GET /api/v1/workflows/:name
func GetWorkflowHandler(c *gin.Context) {}

// GetWorkflowListHandler GET /api/v1/workflows
func GetWorkflowListHandler(c *gin.Context) {}

// DeleteWorkflowHandler DELETE /api/v1/workflows/:name
func DeleteWorkflowHandler(c *gin.Context) {}

// UpdateWorkflowHandler PUT /api/v1/workflows/:name
func UpdateWorkflowHandler(c *gin.Context) {}

// TriggerWorkflowHandler POST /api/v1/workflows/:name/trigger
func TriggerWorkflowHandler(c *gin.Context) {}
