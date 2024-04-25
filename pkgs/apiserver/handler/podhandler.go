package handler

import (
	"github.com/gin-gonic/gin"
)

// CreatePodHandler POST /api/v1/namespaces/:namespace/pods
func CreatePodHandler(c *gin.Context) {}

// GetPodHandler GET /api/v1/namespaces/:namespace/pods/:name
func GetPodHandler(c *gin.Context) {}

// GetPodListHandler GET /api/v1/namespaces/:namespace/pods
func GetPodListHandler(c *gin.Context) {}

// DeletePodHandler DELETE /api/v1/namespaces/:namespace/pods/:name
func DeletePodHandler(c *gin.Context) {}

// UpdatePodHandler PUT /api/v1/namespaces/:namespace/pods/:name
func UpdatePodHandler(c *gin.Context) {}

// GetAllPodsHandler GET /api/v1/pods
func GetAllPodsHandler(c *gin.Context) {}
