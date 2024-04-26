package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"net/http"
)

// CreateNodeHandler POST /api/v1/nodes
func CreateNodeHandler(c *gin.Context) {
	var nodeConfig core.Node
	err := c.Bind(&nodeConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, "{\"data\": needs node config type}")
		return
	}
	nodeName := fmt.Sprintf("/registry/nodes/%s", nodeConfig.NodeMetaData.Name)
	var oldNode core.Node
	err = etcdClient.Get(context.Background(), nodeName, &oldNode)
	if err == nil {
		// has old node config
		logger.Infof("has old node: %v", oldNode)
	}
	err = etcdClient.Put(context.Background(), nodeName, nodeConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "{\"data\": etcd cannot store data")
		return
	}
}

// GetNodeHandler GET /api/v1/nodes/:name
func GetNodeHandler(c *gin.Context) {
	name := c.Param("name")
	var nodeConfig core.Node
	err := etcdClient.Get(context.Background(), fmt.Sprintf("/registry/nodes/%s", name), &nodeConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "{\"data\": etcd cannot read data}")
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("{\"data\": %s}", utils.CreateJson(nodeConfig)))
}

// GetAllNodesHandler GET /api/v1/nodes
func GetAllNodesHandler(c *gin.Context) {
	key := "/registry/nodes/"
	var nodes []core.Node
	err := etcdClient.GetList(context.Background(), key, &nodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "{\"data\": etcd cannot read data}")
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("{\"data\": %s}", utils.CreateJson(nodes)))
}
