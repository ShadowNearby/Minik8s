package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"net/http"
)

func EndpointKeyPrefix(namespace string, name string) string {
	return fmt.Sprintf("/endpoints/%s/object/%s", namespace, name)
}

func EndpointListKeyPrefix(namespace string) string {
	return fmt.Sprintf("/endpoints/%s/object", namespace)
}

// CreateEndpointHandler POST /api/v1/namespaces/:namespace/endpoints
func CreateEndpointHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	endpointConfig := core.Endpoint{}
	if err := c.Bind(&endpointConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := EndpointKeyPrefix(namespace, endpointConfig.MetaData.Name)
	existEndpointConfig := core.Endpoint{}
	if err := storage.Get(key, &existEndpointConfig); err == nil {
		log.Errorf("endpoint %s:%s already exist", namespace, endpointConfig.MetaData.Name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint already exist"})
		return
	}
	if err := storage.Put(key, endpointConfig); err != nil {
		log.Errorf("save endpoint error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetEndpointHandler GET /api/v1/namespaces/:namespace/endpoints/:name
func GetEndpointHandler(c *gin.Context) {
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
	key := EndpointKeyPrefix(namespace, name)
	endpointConfig := core.Endpoint{}
	if err := storage.Get(key, &endpointConfig); err == nil {
		log.Errorf("endpoint %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetEndpointListHandler GET /api/v1/namespaces/:namespace/endpoints
func GetEndpointListHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	key := EndpointListKeyPrefix(namespace)
	var endpointListConfig []core.Endpoint
	if err := storage.RangeGet(key, &endpointListConfig); err == nil {
		log.Errorf("endpoint list %s not found", namespace)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteEndpointHandler DELETE /api/v1/namespaces/:namespace/endpoints/:name
func DeleteEndpointHandler(c *gin.Context) {}

// UpdateEndpointHandler PUT /api/v1/namespaces/:namespace/endpoints/:name
func UpdateEndpointHandler(c *gin.Context) {}
