package controller

import (
	"minik8s/pkgs/controller/replicaset"
	"minik8s/pkgs/controller/service"
)

func Run() {
	var serviceController service.ServiceController
	go StartController(&serviceController)
	var replicaController rsc.ReplicaSetController
	go StartController(&replicaController)
	select {}
}
