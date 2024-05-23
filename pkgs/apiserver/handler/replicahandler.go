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

// CreateReplicaHandler POST /api/v1/namespaces/:namespace/replicas
func CreateReplicaHandler(c *gin.Context) {
	var replica core.ReplicaSet
	//var replicaOld core.ReplicaSet
	err := c.Bind(&replica)
	namespace := c.Param("namespace")
	if err != nil {
		logger.Errorf("bad body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replica set type"})
		return
	}
	namespace = "default"
	rsName := fmt.Sprintf("/replicas/object/%s/%s", namespace, replica.MetaData.Name)
	//if err = storage.Get(rsName, &replicaOld); err == nil {
	//	logger.Errorf("has existed")
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "replica has existed"})
	//	return
	//}
	err = storage.Put(rsName, utils.JsonMarshal(replica))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot put data"})
		return
	}
	// channel
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelReplica, constants.ChannelCreate), utils.JsonMarshal(replica))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetReplicaHandler GET /api/v1/namespaces/:namespace/replicas/:name
func GetReplicaHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	namespace = "default"
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
	namespace = "default"
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
	namespace = "default"
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replica set name in path"})
		return
	}
	var replicaset core.ReplicaSet
	err := storage.Get(fmt.Sprintf("/replicas/object/%s/%s", namespace, name), &replicaset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no such replica set"})
		return
	}
	err = storage.Del(fmt.Sprintf("/replicas/object/%s/%s", namespace, name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del data"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelReplica, constants.ChannelDelete), utils.JsonMarshal(replicaset))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateReplicaHandler PUT /api/v1/namespaces/:namespace/replicas/:name
func UpdateReplicaHandler(c *gin.Context) {
	var newReplica core.ReplicaSet
	var oldReplica core.ReplicaSet
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := c.Bind(&newReplica)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replica set type"})
		return
	}
	namespace = "default" // namespace should be default
	// check
	if name != newReplica.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path info should same as replica info"})
		return
	}
	rsName := fmt.Sprintf("/replicas/object/%s/%s", namespace, name)
	err = storage.Get(rsName, &oldReplica)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not old replicaset"})
		return
	}
	err = storage.Put(rsName, utils.JsonMarshal(newReplica))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot put data"})
		return
	}
	// channel
	replicas := make([]core.ReplicaSet, 2)
	replicas[0] = oldReplica
	replicas[1] = newReplica
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelReplica, constants.ChannelUpdate), utils.JsonMarshal(replicas))
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateReplicaStatusHandler  PUT /api/v1/namespaces/:namespace/replicas/:name/status
func UpdateReplicaStatusHandler(c *gin.Context) {
	// just for update status or owner-reference
	var oldReplica, replica core.ReplicaSet
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := c.Bind(&replica)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect replicaset in body"})
		return
	}
	namespace = "default" // namespace should be default
	if name != replica.MetaData.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expect same name"})
		return
	}
	key := fmt.Sprintf("/replicas/object/%s/%s", namespace, name)
	err = storage.Get(key, &oldReplica)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get old replicaset"})
		return
	}
	// update replica owner-reference
	oldReplica.MetaData.OwnerReference = replica.MetaData.OwnerReference
	// update replica status
	oldReplica.Status = replica.Status
	err = storage.Put(key, oldReplica)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update "})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
