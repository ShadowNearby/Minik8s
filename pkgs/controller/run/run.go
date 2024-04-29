package run

import (
	"minik8s/pkgs/controller"
	"minik8s/pkgs/controller/replicaset"
	"minik8s/pkgs/controller/service"
)

func Run() {
	var serviceController service.ServiceController
	go controller.StartController(&serviceController)
	var replicaController rsc.ReplicaSetController
	go controller.StartController(&replicaController)
	select {}
}
