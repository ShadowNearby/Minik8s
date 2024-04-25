package handler

import "github.com/gin-gonic/gin"

// CreateFunctionHandler POST /api/v1/functions
func CreateFunctionHandler(c *gin.Context) {}

// GetFunctionHandler GET /api/v1/functions/:name
func GetFunctionHandler(c *gin.Context) {}

// DeleteFunctionHandler DELETE /api/v1/functions/:name
func DeleteFunctionHandler(c *gin.Context) {}

// UpdateFunctionHandler PUT /api/v1/functions/:name
func UpdateFunctionHandler(c *gin.Context) {}

// TriggerFunctionHandler POST /api/v1/functions/:name/trigger
func TriggerFunctionHandler(c *gin.Context) {}

// GetAllFunctionsHandler GET /api/v1/functions
func GetAllFunctionsHandler(c *gin.Context) {}
