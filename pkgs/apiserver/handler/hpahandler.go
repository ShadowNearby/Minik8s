package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"
)

// CreateHpaHandler POST /api/v1/namespaces/:namespace/hpa
func CreateHpaHandler(c *gin.Context) {
	namespace := "default"
	var hpa core.HorizontalPodAutoscaler
	err := c.Bind(&hpa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "except hpa in body"})
		return
	}
	key := fmt.Sprintf("/hpa/object/%s/%s", namespace, hpa.MetaData.Name)
	err = storage.Put(key, hpa)
	if err != nil {
		logger.Errorf("put error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelHPA, constants.ChannelCreate), hpa)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// GetHpaHandler GET /api/v1/namespaces/:namespace/hpa/:name
func GetHpaHandler(c *gin.Context) {
	var hpa core.HorizontalPodAutoscaler
	namespace := "default"
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect hpa name"})
		return
	}
	key := fmt.Sprintf("/hpa/object/%s/%s", namespace, name)
	err := storage.Get(key, &hpa)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(hpa)})
}

// GetHpaListHandler GET /api/v1/namespaces/:namespace/hpa
func GetHpaListHandler(c *gin.Context) {
	var hpas []core.HorizontalPodAutoscaler
	namespace := "default"
	key := fmt.Sprintf("/hpa/object/%s", namespace)
	err := storage.RangeGet(key, &hpas)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(hpas)})
}

// DeleteHpaHandler DELETE /api/v1/namespaces/:namespace/hpa/:name
func DeleteHpaHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := "default"
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace and name"})
		return
	}
	hpa := core.HorizontalPodAutoscaler{}
	key := fmt.Sprintf("/hpa/object/%s/%s", namespace, name)
	err := storage.Get(key, &hpa)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	err = storage.Del(key)
	if err != nil {
		logger.Errorf("del error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del data"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelHPA, constants.ChannelDelete), utils.JsonMarshal(hpa))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateHpaHandler PUT /api/v1/namespaces/:namespace/hpa/:name
func UpdateHpaHandler(c *gin.Context) {
	var hpa core.HorizontalPodAutoscaler
	var oldHpa core.HorizontalPodAutoscaler
	namespace := "default"
	name := c.Param("name")
	err := c.Bind(&hpa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	if name != hpa.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace and name should same as path"})
		return
	}
	path := fmt.Sprintf("/hpa/object/%s/%s", namespace, name)
	err = storage.Get(path, &oldHpa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get old pod"})
		return
	}
	err = storage.Put(path, hpa)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	pods := []core.HorizontalPodAutoscaler{oldHpa, hpa}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelHPA, constants.ChannelUpdate), utils.JsonMarshal(pods))
	c.JSON(http.StatusOK, gin.H{})
}

// GetAllHpaHandler GET /api/v1/hpa
func GetAllHpaHandler(c *gin.Context) {
	var hpas []core.HorizontalPodAutoscaler
	namespace := "default"
	key := fmt.Sprintf("/hpa/object/%s", namespace)
	err := storage.RangeGet(key, &hpas)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(hpas)})
}

// UpdateHpaStatusHandler PUT /api/v1/namespaces/:namespace/hpa/:name/status
// there's no side effect of channel in this function, we can update status and owner-reference here
func UpdateHpaStatusHandler(c *gin.Context) {
	var hpa core.HorizontalPodAutoscaler
	var oldHpa core.HorizontalPodAutoscaler
	namespace := "default"
	name := c.Param("name")
	err := c.Bind(&hpa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs pod config type"})
		return
	}
	if name != hpa.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace and name should same as path"})
		return
	}
	path := fmt.Sprintf("/hpa/object/%s/%s", namespace, name)
	err = storage.Get(path, &oldHpa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get old pod"})
		return
	}
	oldHpa.Status = hpa.Status
	oldHpa.MetaData.OwnerReference = hpa.MetaData.OwnerReference
	err = storage.Put(path, oldHpa)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
