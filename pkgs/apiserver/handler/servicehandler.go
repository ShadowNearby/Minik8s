package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"net/http"
)

func ServiceKeyPrefix(namespace string, name string) string {
	return fmt.Sprintf("/services/object/%s/%s", namespace, name)
}

func ServiceListKeyPrefix(namespace string) string {
	return fmt.Sprintf("/services/object/%s", namespace)
}

// CreateServiceHandler POST /api/v1/namespaces/:namespace/services
func CreateServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	serviceConfig := core.Service{}
	if err := c.Bind(&serviceConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := ServiceKeyPrefix(namespace, serviceConfig.MetaData.Name)
	existServiceConfig := core.Service{}
	if err := storage.Get(key, &existServiceConfig); err == nil {
		log.Errorf("service %s:%s already exist", namespace, serviceConfig.MetaData.Name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service already exist"})
		return
	}
	if err := storage.Put(key, serviceConfig); err != nil {
		log.Errorf("save service error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetServiceHandler GET /api/v1/namespaces/:namespace/services/:name
func GetServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need name"})
		return
	}
	key := ServiceKeyPrefix(namespace, name)
	serviceConfig := &core.Service{}
	if err := storage.Get(key, serviceConfig); err == nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetServiceListHandler GET /api/v1/namespaces/:namespace/services
func GetServiceListHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	key := ServiceListKeyPrefix(namespace)
	var serviceListConfig []core.Service
	if err := storage.RangeGet(key, &serviceListConfig); err == nil {
		log.Errorf("service list %s not found", namespace)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteServiceHandler DELETE /api/v1/namespaces/:namespace/services/:name
func DeleteServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need name"})
		return
	}
	key := ServiceKeyPrefix(namespace, name)
	if err := storage.Del(key); err == nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateServiceHandler PUT /api/v1/namespaces/:namespace/services/:name
func UpdateServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need name"})
		return
	}
	serviceConfig := core.Service{}
	if err := c.Bind(&serviceConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := ServiceKeyPrefix(namespace, name)
	if err := storage.Put(key, serviceConfig); err == nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service put error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
