package kubeletcontroller

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/handler"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var KubeletRouter = [...]handler.Route{
	{Path: "/pod/create", Method: "POST", Handler: CreatePodController},
	{Path: "/pod/stop/:namespace/:name", Method: "DELETE", Handler: StopPodController},
	{Path: "/pod/status/:namespace/:name", Method: "GET", Handler: InspectPodController}, // running status
	{Path: "/metrics/:namespace/:name", Method: "GET", Handler: PodMetricController},     // for auto-scaling
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
	if err != nil {
		StopPod(pod)
		c.JSON(http.StatusInternalServerError, "cannot create pod")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(pod.Status)})
}

func StopPodController(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	pod := runtime.KubeletInstance.GetPodConfig(name, namespace)
	err := StopPod(pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func InspectPodController(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	pConfig := runtime.KubeletInstance.GetPodConfig(name, namespace)
	inspect := InspectPod(&pConfig, runtime.ExecProbe)
	c.JSON(http.StatusOK, gin.H{"data": inspect})
}

func PodMetricController(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	pod := runtime.KubeletInstance.GetPodConfig(name, namespace)
	metric := PodMetrics(pod)
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(metric)})
}

func NodeMetricsController(c *gin.Context) {
	metrics := NodeMetrics()
	text := utils.JsonMarshal(metrics)
	c.JSON(http.StatusOK, gin.H{"data": text})
}
