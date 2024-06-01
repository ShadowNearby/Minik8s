package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ServiceKeyPrefix(namespace string, name string) string {
	return fmt.Sprintf("/services/object/%s/%s", namespace, name)
}

func ServiceListKeyPrefix(namespace string) string {
	return fmt.Sprintf("/services/object/%s", namespace)
}

// CreateServiceHandler POST /api/v1/namespaces/:namespace/services
func CreateServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	serviceConfig := core.Service{}
	if err := c.Bind(&serviceConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := ServiceKeyPrefix(namespace, serviceConfig.MetaData.Name)
	if err := storage.Put(key, serviceConfig); err != nil {
		log.Errorf("save service error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelService, constants.ChannelCreate), utils.JsonMarshal(serviceConfig))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetServiceHandler GET /api/v1/namespaces/:namespace/services/:name
func GetServiceHandler(c *gin.Context) {
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
	key := ServiceKeyPrefix(namespace, name)
	serviceConfig := &core.Service{}
	if err := storage.Get(key, serviceConfig); err != nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(serviceConfig)})
}

// GetServiceListHandler GET /api/v1/namespaces/:namespace/services
func GetServiceListHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	key := ServiceListKeyPrefix(namespace)
	var serviceListConfig []core.Service
	if err := storage.RangeGet(key, &serviceListConfig); err != nil {
		log.Errorf("service list %s not found", namespace)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(serviceListConfig)})
}

// DeleteServiceHandler DELETE /api/v1/namespaces/:namespace/services/:name
func DeleteServiceHandler(c *gin.Context) {
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
	key := ServiceKeyPrefix(namespace, name)
	serviceConfig := &core.Service{}
	if err := storage.Get(key, serviceConfig); err != nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	if err := storage.Del(key); err != nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelService, constants.ChannelDelete), utils.JsonMarshal(serviceConfig))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateServiceHandler PUT /api/v1/namespaces/:namespace/services/:name
func UpdateServiceHandler(c *gin.Context) {
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
	serviceConfig := &core.Service{}
	if err := c.Bind(serviceConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := ServiceKeyPrefix(namespace, name)
	preServiceConfig := &core.Service{}
	if err := storage.Get(key, preServiceConfig); err != nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not found"})
		return
	}
	if err := storage.Put(key, serviceConfig); err != nil {
		log.Errorf("service %s:%s not found", namespace, name)
		c.JSON(http.StatusBadRequest, gin.H{"error": "service put error"})
		return
	}
	services := []core.Service{*preServiceConfig, *serviceConfig}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelService, constants.ChannelUpdate), utils.JsonMarshal(services))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

const ServiceIPKey = "/services/ip"

var UsedMapKey = fmt.Sprintf("%s/used", ServiceIPKey)

var ServiceIPMapKey = fmt.Sprintf("%s/map", ServiceIPKey)

func GetServiceClusterIPHandler(c *gin.Context) {
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
	usedMap := map[string]bool{}
	serviceIPMap := map[string]string{}
	serviceKey := fmt.Sprintf("%s:%s", namespace, name)
	err1 := storage.Get(UsedMapKey, &usedMap)
	err2 := storage.Get(ServiceIPMapKey, &serviceIPMap)
	var newIP string
	if err1 != nil || err2 != nil {
		newIP = utils.GenerateNewClusterIP()
	} else {
		for {
			newIP = utils.GenerateNewClusterIP()
			if exist, ok := usedMap[newIP]; !ok || !exist {
				break
			}
		}
	}
	oldIP, ok := serviceIPMap[serviceKey]
	if ok {
		c.JSON(http.StatusOK, gin.H{"data": oldIP})
		return
	}
	usedMap[newIP] = true
	serviceIPMap[serviceKey] = newIP
	err := storage.Put(UsedMapKey, usedMap)
	if err != nil {
		log.Errorf("error in put UsedMap")
		c.JSON(http.StatusOK, gin.H{"error": "error in put UsedMap"})
	}
	err = storage.Put(ServiceIPMapKey, serviceIPMap)
	if err != nil {
		log.Errorf("error in put ServiceIPMap")
		c.JSON(http.StatusOK, gin.H{"error": "error in put ServiceIPMap"})
	}
	c.JSON(http.StatusOK, gin.H{"data": newIP})
}

func DeleteServiceClusterIPHandler(c *gin.Context) {
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
	usedMap := map[string]bool{}
	serviceIPMap := map[string]string{}
	serviceKey := fmt.Sprintf("%s:%s", namespace, name)
	err1 := storage.Get(UsedMapKey, &usedMap)
	err2 := storage.Get(ServiceIPMapKey, &serviceIPMap)
	if err1 != nil || err2 != nil {
		log.Errorf("error in get ServiceIPMap or UsedMapKey")
		c.JSON(http.StatusOK, gin.H{"error": "error in get ServiceIPMap or UsedMapKey"})
	}
	ip, ok := serviceIPMap[serviceKey]
	if !ok {
		log.Errorf("error in find ip for service %s", serviceKey)
		c.JSON(http.StatusOK, gin.H{"error": "error in find ip for service"})
	}
	delete(serviceIPMap, serviceKey)
	delete(usedMap, ip)
	err := storage.Put(UsedMapKey, usedMap)
	if err != nil {
		log.Errorf("error in put UsedMap")
		c.JSON(http.StatusOK, gin.H{"error": "error in put UsedMap"})
	}
	err = storage.Put(ServiceIPMapKey, serviceIPMap)
	if err != nil {
		log.Errorf("error in put ServiceIPMap")
		c.JSON(http.StatusOK, gin.H{"error": "error in put ServiceIPMap"})
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
