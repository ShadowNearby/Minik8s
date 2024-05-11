package test

import (
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/apiserver/storage"
)

func ServerRun() {
	server := server.CreateAPIServer(storage.DefaultEndpoints)
	server.Run(fmt.Sprintf("%s:%s", config.LocalServerIp, config.ApiServerPort))
}
