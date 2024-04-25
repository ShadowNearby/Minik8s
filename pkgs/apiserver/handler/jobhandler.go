package handler

import "github.com/gin-gonic/gin"

// CreateJobHandler POST /api/v1/namespaces/:namespace/jobs
func CreateJobHandler(c *gin.Context) {}

// GetJobHandler GET /api/v1/namespaces/:namespace/jobs/:name
func GetJobHandler(c *gin.Context) {}

// GetJobListHandler GET /api/v1/namespaces/:namespace/jobs
func GetJobListHandler(c *gin.Context) {}

// DeleteJobHandler DELETE /api/v1/namespaces/:namespace/jobs/:name
func DeleteJobHandler(c *gin.Context) {}

// UpdateJobHandler PUT /api/v1/namespaces/:namespace/jobs/:name
func UpdateJobHandler(c *gin.Context) {}
