package apiserver

import (
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/controller"
	scheduler "minik8s/pkgs/controller/scheduler"
	"minik8s/pkgs/controller/service"
)

func Run() {
	server := server.CreateAPIServer(config.DefaultEtcdEndpoints)
	var serviceController service.ServiceController
	go controller.StartController(&serviceController)
	var schedulerController scheduler.Scheduler
	go schedulerController.Run(config.PolicyCPU)
	server.Run(fmt.Sprintf("%s:%s", config.LocalServerIp, config.ApiServerPort))
}
