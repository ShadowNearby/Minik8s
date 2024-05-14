package test

import (
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/apiserver/server"
)

func ServerRun() {
	server := server.CreateAPIServer(config.DefaultEtcdEndpoints)
	server.Run(fmt.Sprintf("%s:%s", config.LocalServerIp, config.ApiServerPort))
}
