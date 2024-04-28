package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"net/http"
)

// CreateNodeHandler POST /api/v1/nodes
func CreateNodeHandler(c *gin.Context) {
	var nodeConfig core.Node
	err := c.Bind(&nodeConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs node config type"})
		return
	}
	nodeName := fmt.Sprintf("/nodes/object/%s", nodeConfig.NodeMetaData.Name)
	err = storage.Put(nodeName, nodeConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetNodeHandler GET /api/v1/nodes/:name
func GetNodeHandler(c *gin.Context) {
	var node core.Node
	name := c.Param("name")
	key := fmt.Sprintf("/nodes/object/%s", name)
	err := storage.Get(key, &node)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot find resource"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(node)})
}

// GetAllNodesHandler GET /api/v1/nodes
func GetAllNodesHandler(c *gin.Context) {
	key := "/nodes/object"
	var nodes []core.Node
	err := storage.RangeGet(key, &nodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(nodes)})
}
