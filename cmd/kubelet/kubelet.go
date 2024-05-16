package main

import (
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet"

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
	kubelet.Run(config, "127.0.0.1:10250")
}
