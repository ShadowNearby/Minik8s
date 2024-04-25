package handler

import "github.com/gin-gonic/gin"

// CreateDNSHandler POST /api/v1/namespaces/:namespace/dns
func CreateDNSHandler(c *gin.Context) {}

// GetDNSHandler GET /api/v1/namespaces/:namespace/dns/:name
func GetDNSHandler(c *gin.Context) {}

// GetDNSListHandler GET /api/v1/namespaces/:namespace/dns
func GetDNSListHandler(c *gin.Context) {}

// DeleteDNSHandler DELETE /api/v1/namespaces/:namespace/dns/:name
func DeleteDNSHandler(c *gin.Context) {}

// UpdateDNSHandler PUT /api/v1/namespaces/:namespace/dns/:name
func UpdateDNSHandler(c *gin.Context) {}
