package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"net/http"
)

var KubeletRouter = [...]core.Route{
	{Path: "/api/v1/kubelet/pod/create", Method: "POST", Handler: CreatePodController},
	{Path: "/api/vi/kubelet/pod/stop", Method: "POST", Handler: StopPodController},
	{Path: "/:namespace/:podName", Method: "GET", Handler: InspectPodController},
	{Path: "/metrics", Method: "GET", Handler: NodeMetricsController},
}

func CreatePodController(c *gin.Context) {
	json := make(map[string]interface{})
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	podConfigText := json["podConfig"]
	if podConfigText == nil || podConfigText == "" {
		c.JSON(http.StatusBadRequest, "expect podConfig")
		return
	}
	podConfig := podConfigText.(core.Pod)
	err = CreatePod(&podConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "cannot create pod")
		return
	}
	c.JSON(http.StatusOK, "")
}

func StopPodController(c *gin.Context) {
	json := make(map[string]interface{})
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	podConfigText := json["podConfig"]
	if podConfigText == nil || podConfigText == "" {
		c.JSON(http.StatusBadRequest, "expect podConfig")
		return
	}
	podConfig := podConfigText.(core.Pod)
	err = StopPod(podConfig)
	if err != nil {
		return
	}
}

func InspectPodController(c *gin.Context) {
	name := c.Param("podName")
	namespace := c.Param("namespace")
	pConfig := runtime.KubeletInstance.GetPodConfig(name, namespace)
	inspect := InspectPod(&pConfig, runtime.ExecProbe)
	c.JSON(http.StatusOK, fmt.Sprintf("{\"status\": %s}", inspect))
}

func NodeMetricsController(c *gin.Context) {
	metrics := NodeMetrics()
	text := utils.JsonMarshal(metrics)
	c.JSON(http.StatusOK, fmt.Sprintf("{\"metrics\":%s}", text))
}
