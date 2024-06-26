package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func EndpointKeyPrefix(namespace string, name string) string {
	return fmt.Sprintf("/endpoints/object/%s/%s", namespace, name)
}

func EndpointListKeyPrefix(namespace string) string {
	return fmt.Sprintf("/endpoints/object/%s", namespace)
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
	if err := storage.Put(key, endpointConfig); err != nil {
		log.Errorf("save endpoint error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
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
	if err := storage.Get(key, &endpointConfig); err != nil {
		log.Errorf("endpoint %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(endpointConfig)})
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
	if err := storage.RangeGet(key, &endpointListConfig); err != nil {
		log.Errorf("endpoint list %s not found", namespace)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(endpointListConfig)})
}

// DeleteEndpointHandler DELETE /api/v1/namespaces/:namespace/endpoints/:name
func DeleteEndpointHandler(c *gin.Context) {
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
	if err := storage.Del(key); err != nil {
		log.Errorf("endpoint %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateEndpointHandler PUT /api/v1/namespaces/:namespace/endpoints/:name
func UpdateEndpointHandler(c *gin.Context) {
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
	endpointConfig := core.Endpoint{}
	if err := c.Bind(&endpointConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := EndpointKeyPrefix(namespace, name)
	if err := storage.Put(key, endpointConfig); err != nil {
		log.Errorf("endpoint %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
