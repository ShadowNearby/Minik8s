package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

// 1. update -> no need to reschedule
// 2. create, replace -> need reschedule
// 3. delete -> tell kubelet to stop pod

// CreatePodHandler POST /api/v1/namespaces/:namespace/pods
func CreatePodHandler(c *gin.Context) {
	var pod core.Pod
	err := c.Bind(&pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	if pod.MetaData.Namespace == "" {
		pod.MetaData.Namespace = "default"
	}
	key := fmt.Sprintf("/pods/object/%s/%s", pod.MetaData.Namespace, pod.MetaData.Name)
	err = storage.Put(key, pod)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelPod, constants.ChannelCreate), pod)
	pods := []core.Pod{{}, pod}
	storage.RedisInstance.PublishMessage(constants.ChannelPodSchedule, utils.JsonMarshal(pods))
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
	var pod core.Pod
	err := storage.Get(podName, &pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no such pod"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelPod, constants.ChannelDelete), pod)
	err = storage.Del(podName)
	if err != nil {
		logger.Errorf("del error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdatePodHandler PUT /api/v1/namespaces/:namespace/pods/:name
func UpdatePodHandler(c *gin.Context) {
	var pod core.Pod
	var oldPod core.Pod
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := c.Bind(&pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	if namespace != pod.MetaData.Namespace || name != pod.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace and name should same as path"})
		return
	}
	path := fmt.Sprintf("/pods/object/%s/%s", pod.MetaData.Namespace, pod.MetaData.Name)
	err = storage.Get(path, &oldPod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get old pod"})
		return
	}
	err = storage.Put(path, pod)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	pods := []core.Pod{oldPod, pod}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelPod, constants.ChannelUpdate), utils.JsonMarshal(pods))
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

// ReplacePodHandler POST /api/v1/namespaces/:namespace/pods/:name
func ReplacePodHandler(c *gin.Context) {
	var pod core.Pod
	var oldPod core.Pod
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := c.Bind(&pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	if namespace != pod.MetaData.Namespace || name != pod.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace and name should same as path"})
		return
	}
	path := fmt.Sprintf("/pods/object/%s/%s", pod.MetaData.Namespace, pod.MetaData.Name)
	err = storage.Get(path, &oldPod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get old pod"})
		return
	}
	// 1. hostip
	// 2. labels
	// need to reschedule
	pods := []core.Pod{oldPod, pod}
	storage.RedisInstance.PublishMessage(constants.ChannelPodSchedule, utils.JsonMarshal(pods))
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
