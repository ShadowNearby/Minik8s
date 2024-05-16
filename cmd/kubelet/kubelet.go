package main

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet"
	"minik8s/utils"

	logger "github.com/sirupsen/logrus"
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{DisableTimestamp: true})
	logger.SetReportCaller(true)
	var config = core.KubeletConfig{
		MasterIP:   config.LocalServerIp,
		MasterPort: config.ApiServerPort,
		Labels: map[string]string{
			"test": "haha",
		},
	}
	kubelet.Run(config, fmt.Sprintf("%s:%d", utils.GetIP(), 10250))
}
