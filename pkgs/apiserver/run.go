package apiserver

import (
	"fmt"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/config"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/controller/service"
)

func Run() {
	server := server.CreateAPIServer(storage.DefaultEndpoints)
	var serviceController service.ServiceController
	go controller.StartController(&serviceController)
	server.Run(fmt.Sprintf("%s:%s", config.LocalServerIp, config.ApiServerPort))
}
