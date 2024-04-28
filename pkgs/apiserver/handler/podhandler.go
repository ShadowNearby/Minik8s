package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"net/http"
)

// CreatePodHandler POST /api/v1/namespaces/:namespace/pods
func CreatePodHandler(c *gin.Context) {
	var podConfig core.Pod
	err := c.Bind(&podConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	podName := fmt.Sprintf("/pods/object/%s/%s", podConfig.MetaData.NameSpace, podConfig.MetaData.Name)
	err = storage.Put(podName, podConfig)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	storage.RedisInstance.PublishMessage(storage.ChannelNewPod, podName)
	c.JSON(http.StatusOK, gin.H{})
}

// GetPodHandler GET /api/v1/namespaces/:namespace/pods/:name
func GetPodHandler(c *gin.Context) {
	var podConfig core.Pod
	name := c.Param("name")
	namespace := c.Param("namespace")
	if len(name) == 0 || len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace and name"})
		return
	}
	podName := fmt.Sprintf("/pods/object/%s/%s", namespace, name)
	err := storage.Get(podName, &podConfig)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(podConfig)})
}

// GetPodListHandler GET /api/v1/namespaces/:namespace/pods
func GetPodListHandler(c *gin.Context) {
	var podConfigs []core.Pod
	namespace := c.Param("namespace")
	if len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace"})
		return
	}
	podName := fmt.Sprintf("/pods/object/%s", namespace)
	err := storage.RangeGet(podName, &podConfigs)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(podConfigs)})
}

// DeletePodHandler DELETE /api/v1/namespaces/:namespace/pods/:name
func DeletePodHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if len(name) == 0 || len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace and name"})
		return
	}
	podName := fmt.Sprintf("/pods/object/%s/%s", namespace, name)
	err := storage.Del(podName)
	if err != nil {
		logger.Errorf("del error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdatePodHandler PUT /api/v1/namespaces/:namespace/pods/:name
func UpdatePodHandler(c *gin.Context) {
	var podConfig core.Pod
	err := c.Bind(&podConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	podName := fmt.Sprintf("/pods/object/%s/%s", podConfig.MetaData.NameSpace, podConfig.MetaData.Name)
	err = storage.Put(podName, podConfig)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	storage.RedisInstance.PublishMessage(storage.ChannelUpdatePod, podName)
	c.JSON(http.StatusOK, gin.H{})
}

// GetAllPodsHandler GET /api/v1/pods
func GetAllPodsHandler(c *gin.Context) {
	var podConfigs []core.Pod
	err := storage.RangeGet("/pods", &podConfigs)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(podConfigs)})
}
