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

// CreateJobHandler POST /api/v1/namespaces/:namespace/jobs
func CreateJobHandler(c *gin.Context) {
	var job core.Job
	err := c.Bind(&job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check the parameters
	if job.MetaData.Name == "" {
		logger.Errorf("Job name is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "function is empty"})
		return
	}
	if job.MetaData.Namespace == "" {
		job.MetaData.Namespace = "default"
	}
	key := fmt.Sprintf("/jobs/%s/%s", job.MetaData.Namespace, job.MetaData.Name)
	job.Status = core.PodStatus{Phase: core.PodPhasePending}
	err = storage.Put(key, job)
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelJob, constants.ChannelCreate), utils.JsonMarshal(job))
	logger.Info("[create jod successfully]")
	c.JSON(http.StatusOK, gin.H{"data": "create job success"})
}

// GetJobHandler GET /api/v1/namespaces/:namespace/jobs/:name
func GetJobHandler(c *gin.Context) {
	var job core.Job
	key := fmt.Sprintf("/jobs/%s/%s", job.MetaData.Namespace, job.MetaData.Name)
	err := storage.Get(key, &job)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get function"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(job)})
}

// GetJobListHandler GET /api/v1/namespaces/:namespace/jobs
func GetJobListHandler(c *gin.Context) {
	var jobConfigs []core.Job
	namespace := c.Param("namespace")
	if len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace"})
		return
	}
	jobName := fmt.Sprintf("/jobs/%s", namespace)
	err := storage.RangeGet(jobName, &jobConfigs)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(jobConfigs)})
}

// DeleteJobHandler DELETE /api/v1/namespaces/:namespace/jobs/:name
func DeleteJobHandler(c *gin.Context) {
	var job core.Job
	err := c.Bind(&job)
	key := fmt.Sprintf("/jobs/%s/%s", job.MetaData.Namespace, job.MetaData.Name)
	err = storage.Del(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete function"})
		return
	}
	storage.RedisInstance.PublishMessage(
		constants.GenerateChannelName(constants.ChannelJob, constants.ChannelDelete), job.MetaData.Name)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// UpdateJobHandler PUT /api/v1/namespaces/:namespace/jobs/:name
func UpdateJobHandler(c *gin.Context) {

	var oldJob, newJob core.Job
	err := c.Bind(&newJob)
	key := fmt.Sprintf("/jobs/%s/%s", newJob.MetaData.Namespace, newJob.MetaData.Name)
	err = storage.Get(key, &oldJob)
	err = storage.Put(key, newJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot set object"})
		return
	}
	jobs := []core.Job{oldJob, newJob}
	storage.RedisInstance.PublishMessage(constants.GenerateChannelName(constants.ChannelJob, constants.ChannelUpdate), utils.JsonMarshal(jobs))
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
