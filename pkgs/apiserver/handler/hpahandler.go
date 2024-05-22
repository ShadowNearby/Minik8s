package handler

import "github.com/gin-gonic/gin"

// CreateHpaHandler POST /api/v1/namespaces/:namespace/hpa
func CreateHpaHandler(c *gin.Context) {}

// GetHpaHandler GET /api/v1/namespaces/:namespace/hpa/:name
func GetHpaHandler(c *gin.Context) {}

// GetHpaListHandler GET /api/v1/namespaces/:namespace/hpa
func GetHpaListHandler(c *gin.Context) {}

// DeleteHpaHandler DELETE /api/v1/namespaces/:namespace/hpa/:name
func DeleteHpaHandler(c *gin.Context) {}

// UpdateHpaHandler PUT /api/v1/namespaces/:namespace/hpa/:name
func UpdateHpaHandler(c *gin.Context) {}

// GetAllHpaHandler GET /api/v1/hpa
func GetAllHpaHandler(c *gin.Context) {}

// UpdateHpaStatusHandler PUT /api/v1/namespaces/:namespace/hpa/:name/status
func UpdateHpaStatusHandler(c *gin.Context) {}
