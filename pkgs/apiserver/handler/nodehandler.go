package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
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
	nodeConfig.Status.Phase = core.NodeReady
	nodeConfig.Status.LastHeartbeat = time.Now()
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

// DeleteNodeHandler DELETE /api/v1/nodes
func DeleteNodeHandler(c *gin.Context) {
	name := c.Param("name")
	key := fmt.Sprintf("/nodes/object/%s", name)
	var node core.Node
	// TODO: all pods on the node should be scheduled
	err := storage.Get(key, &node)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot get node"})
		return
	}
	err = storage.Del(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot delete node"})
		return
	}
	var pods []core.Pod
	err = storage.RangeGet(fmt.Sprintf("/pods/object/"), &pods)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	for _, pod := range pods {
		logger.Infof("podip:%s, nodeip: %s", pod.Status.HostIP, node.Spec.NodeIP)
		if pod.Status.HostIP == node.Spec.NodeIP {
			storage.RedisInstance.PublishMessage(constants.ChannelPodSchedule, pod)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// UpdateNodeHandler PUT /api/v1/nodes/:name
func UpdateNodeHandler(c *gin.Context) {
	name := c.Param("name")
	key := fmt.Sprintf("/nodes/object/%s", name)
	var node core.Node
	err := c.Bind(&node)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect node in body"})
		return
	}
	err = storage.Get(key, &node)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "does not exist before"})
		return
	}
	err = storage.Put(key, node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
