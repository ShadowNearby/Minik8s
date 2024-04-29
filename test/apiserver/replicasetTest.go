package apiserver

import (
	"fmt"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/config"
)

func ServerRun() {
	server := server.CreateAPIServer(storage.DefaultEndpoints)
	server.Run(fmt.Sprintf("%s:%s", config.LocalServerIp, config.ApiServerPort))
}
