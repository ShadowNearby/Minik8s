package handler

import "github.com/gin-gonic/gin"

// CreateEndpointHandler POST /api/v1/namespaces/:namespace/endpoints
func CreateEndpointHandler(c *gin.Context) {}

// GetEndpointHandler GET /api/v1/namespaces/:namespace/endpoints/:name
func GetEndpointHandler(c *gin.Context) {}

// GetEndpointListHandler GET /api/v1/namespaces/:namespace/endpoints
func GetEndpointListHandler(c *gin.Context) {}

// DeleteEndpointHandler DELETE /api/v1/namespaces/:namespace/endpoints/:name
func DeleteEndpointHandler(c *gin.Context) {}

// UpdateEndpointHandler PUT /api/v1/namespaces/:namespace/endpoints/:name
func UpdateEndpointHandler(c *gin.Context) {}
