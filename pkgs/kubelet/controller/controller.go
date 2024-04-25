package controller

import (
	"github.com/gin-gonic/gin"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"net/http"
)

var KubeletRouter = [...]core.Route{
	{Path: "/api/v1/kubelet/pod/create", Method: constants.MethodPost, Handler: CreatePodController},
	{Path: "/api/vi/kubelet/pod/stop", Method: constants.MethodPost, Handler: StopPodController},
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

}
