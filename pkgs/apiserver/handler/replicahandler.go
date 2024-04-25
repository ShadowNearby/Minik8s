package handler

import "github.com/gin-gonic/gin"

// CreateReplicaHandler POST /api/v1/namespaces/:namespace/replicas
func CreateReplicaHandler(c *gin.Context) {}

// GetReplicaHandler GET /api/v1/namespaces/:namespace/replicas/:name
func GetReplicaHandler(c *gin.Context) {}

// GetReplicaListHandler GET /api/v1/namespaces/:namespace/replicas
func GetReplicaListHandler(c *gin.Context) {}

// DeleteReplicaHandler DELETE /api/v1/namespaces/:namespace/replicas/:name
func DeleteReplicaHandler(c *gin.Context) {}

// UpdateReplicaHandler PUT /api/v1/namespaces/:namespace/replicas/:name
func UpdateReplicaHandler(c *gin.Context) {}
