package kubeletcontroller

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/handler"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"net/http"
)

var KubeletRouter = [...]handler.Route{
	{Path: "/pod/create", Method: "POST", Handler: CreatePodController},
	{Path: "/pod/stop", Method: "POST", Handler: StopPodController},
	{Path: "/:namespace/:podName", Method: "GET", Handler: InspectPodController},
	{Path: "/metrics", Method: "GET", Handler: NodeMetricsController},
}

func CreatePodController(c *gin.Context) {
	var pod core.Pod
	err := c.BindJSON(&pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong data type, expect pod type"})
		return
	}
	err = CreatePod(&pod)
	logger.Info(pod.Status.ContainersStatus)
	if err != nil {
		StopPod(pod)
		c.JSON(http.StatusInternalServerError, "cannot create pod")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func StopPodController(c *gin.Context) {
	var pod core.Pod
	err := c.BindJSON(&pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	pod = runtime.KubeletInstance.GetPodConfig(pod.MetaData.Name, pod.MetaData.Namespace)
	err = StopPod(pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func InspectPodController(c *gin.Context) {
	name := c.Param("podName")
	namespace := c.Param("namespace")
	pConfig := runtime.KubeletInstance.GetPodConfig(name, namespace)
	inspect := InspectPod(&pConfig, runtime.ExecProbe)
	c.JSON(http.StatusOK, gin.H{"data": inspect})
}

func NodeMetricsController(c *gin.Context) {
	metrics := NodeMetrics()
	text := utils.JsonMarshal(metrics)
	c.JSON(http.StatusOK, gin.H{"data": text})
}
