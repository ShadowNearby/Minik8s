package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"net/http"
)

// CreateReplicaHandler POST /api/v1/namespaces/:namespace/replicas
func CreateReplicaHandler(c *gin.Context) {
	var replica core.ReplicaSet
	err := c.Bind(&replica)
	namespace := c.Param("namespace")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replica set type"})
		return
	}
	rsName := fmt.Sprintf("/replicas/object/%s/%s", namespace, replica.MetaData.Name)
	err = storage.Put(rsName, utils.JsonMarshal(replica))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot put data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetReplicaHandler GET /api/v1/namespaces/:namespace/replicas/:name
func GetReplicaHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" {
		namespace = "default"
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect name in path"})
		return
	}
	var replica core.ReplicaSet
	err := storage.Get(fmt.Sprintf("/replicas/object/%s/%s", namespace, name), &replica)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, utils.JsonMarshal(replica))
}

// GetReplicaListHandler GET /api/v1/namespaces/:namespace/replicas
func GetReplicaListHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	var replicas []core.ReplicaSet
	err := storage.RangeGet(fmt.Sprintf("/replicas/object/%s", namespace), &replicas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, utils.JsonMarshal(replicas))
}

// DeleteReplicaHandler DELETE /api/v1/namespaces/:namespace/replicas/:name
func DeleteReplicaHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" {
		namespace = "default"
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replica set name in path"})
		return
	}
	err := storage.Del(fmt.Sprintf("/replicas/object/%s/%s", namespace, name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateReplicaHandler PUT /api/v1/namespaces/:namespace/replicas/:name
func UpdateReplicaHandler(c *gin.Context) {
	var replica core.ReplicaSet
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := c.Bind(&replica)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replica set type"})
		return
	}
	if namespace == "" {
		namespace = "default"
	}
	// check
	if namespace != replica.MetaData.NameSpace || name != replica.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path info should same as replica info"})
		return
	}
	rsName := fmt.Sprintf("/replicas/object/%s/%s", namespace, name)
	err = storage.Put(rsName, utils.JsonMarshal(replica))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot put data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
