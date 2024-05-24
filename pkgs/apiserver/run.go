package apiserver

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/controller/autoscaler"
	"minik8s/pkgs/controller/podcontroller"
	rsc "minik8s/pkgs/controller/replicaset"
	scheduler "minik8s/pkgs/controller/scheduler"
	"minik8s/pkgs/controller/service"
	"minik8s/utils"
)

func Run() {
	config.ClusterMasterIP = utils.GetIP()
	server := server.CreateAPIServer(config.DefaultEtcdEndpoints)
	var serviceController service.ServiceController
	go controller.StartController(&serviceController)
	var schedulerController scheduler.Scheduler
	go schedulerController.Run(config.PolicyCPU)
	var podcontroller podcontroller.PodController
	go controller.StartController(&podcontroller)
	var replicaSet rsc.ReplicaSetController
	go controller.StartController(&replicaSet)
	var hpa autoscaler.HPAController
	go controller.StartController(&hpa)
	// start hpa background work
	go hpa.StartBackground()

	utils.GenerateNginxFile([]core.DNSRecord{})
	utils.StartNginx()

	server.Run(fmt.Sprintf("%s:%s", config.LocalServerIp, config.ApiServerPort))
}
