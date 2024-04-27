package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"net/http"
)

// CreateServiceHandler POST /api/v1/namespaces/:namespace/services
func CreateServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	var serviceConfig core.Service
	if err := c.Bind(&serviceConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := fmt.Sprintf("/services/object/%s", serviceConfig.MetaData.Name)
	err := storage.Put(key, serviceConfig)
	if err != nil {
		log.Errorf("save service error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetServiceHandler GET /api/v1/namespaces/:namespace/services/:name
func GetServiceHandler(c *gin.Context) {}

// GetServiceListHandler GET /api/v1/namespaces/:namespace/services
func GetServiceListHandler(c *gin.Context) {}

// DeleteServiceHandler DELETE /api/v1/namespaces/:namespace/services/:name
func DeleteServiceHandler(c *gin.Context) {}

// UpdateServiceHandler PUT /api/v1/namespaces/:namespace/services/:name
func UpdateServiceHandler(c *gin.Context) {}
