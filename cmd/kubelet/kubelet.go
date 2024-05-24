package main

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet"
	"minik8s/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})
	var config = core.KubeletConfig{
		MasterIP:   config.LocalServerIp,
		MasterPort: config.ApiServerPort,
		Labels: map[string]string{
			"test": "haha",
			"app":  "nginx",
		},
	}
	kubelet.Run(config, fmt.Sprintf("%s:%d", utils.GetIP(), 10250))
}
