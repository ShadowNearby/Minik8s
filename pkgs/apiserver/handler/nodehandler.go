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
		c.JSON(http.StatusBadRequest, httpData("needs node config type"))
		return
	}
	nodeName := fmt.Sprintf("/nodes/object/%s", nodeConfig.NodeMetaData.Name)
	err = storage.Put(nodeName, nodeConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpData("cannot store data"))
	}
	c.JSON(http.StatusOK, "")
}

// GetNodeHandler GET /api/v1/nodes/:name
func GetNodeHandler(c *gin.Context) {
	var node core.Node
	name := c.Param("name")
	key := fmt.Sprintf("/nodes/object/%s", name)
	err := storage.Get(key, node)
	if err != nil {
		c.JSON(http.StatusNotFound, httpData("cannot find resource"))
		return
	}
	c.JSON(http.StatusOK, httpData(utils.JsonMarshal(node)))
}

// GetAllNodesHandler GET /api/v1/nodes
func GetAllNodesHandler(c *gin.Context) {
	key := "/nodes/object"
	var nodes []core.Node
	err := storage.RangeGet(key, &nodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpData("cannot read data"))
		return
	}
	c.JSON(http.StatusOK, httpData(utils.JsonMarshal(nodes)))
}
