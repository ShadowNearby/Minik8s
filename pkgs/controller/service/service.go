package service

import (
	"encoding/json"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/kubeproxy"
	"minik8s/utils"

	log "github.com/sirupsen/logrus"
)

const TotalIP = (1 << 8)

var UsedIP = [TotalIP]bool{}

const IPPrefix = "10.10.0."

type ServiceController struct{}

func (sc *ServiceController) GetChannel() string {
	return constants.ChannelService
}

func (sc *ServiceController) HandleCreate(message string) error {
	service := &core.Service{}
	err := json.Unmarshal([]byte(message), service)
	if err != nil {
		log.Errorf("unmarshal service error: %s", err.Error())
		return err
	}
	var clusterIP string
	if service.Spec.ClusterIP == "" {
		clusterIP := FindUnusedIP()
		service.Spec.ClusterIP = clusterIP
		utils.SetObject(core.ObjService, service.MetaData.Namespace, service.MetaData.Name, service)
	} else {
		clusterIP = service.Spec.ClusterIP
	}

	// creaete service and alloc ip
	if service.Spec.Type == core.ServiceTypeClusterIP {
		for _, port := range service.Spec.Ports {
			kubeproxy.CreateService(clusterIP, port.Port)
		}
	} else if service.Spec.Type == core.ServiceTypeNodePort {
		NodeIP := utils.GetIP()
		for _, port := range service.Spec.Ports {
			kubeproxy.CreateService(NodeIP, port.NodePort)
		}
	}

	err = CreateEndpointObject(service)
	if err != nil {
		log.Errorf("error in CreateEndpointObject")
		return err
	}
	return nil
}

func (sc *ServiceController) HandleUpdate(message string) error {
	services := []core.Service{}
	err := json.Unmarshal([]byte(message), &services)
	if err != nil {
		log.Errorf("unmarshal service error: %s", err.Error())
		return err
	}
	if len(services) != 2 {
		return fmt.Errorf("service update error")
	}
	preService := &services[0]
	service := &services[1]
	previousSelector := preService.Spec.Selector
	if utils.MatchLabel(previousSelector.MatchLabels, service.Spec.Selector.MatchLabels) {
		return nil
	}

	err = DeleteEndpointObject(service)
	if err != nil {
		log.Errorf("error in UpdateEndpointObject")
		return err
	}
	err = CreateEndpointObject(service)
	if err != nil {
		log.Errorf("error in UpdateEndpointObject")
		return err
	}

	return nil
}

func (sc *ServiceController) HandleDelete(message string) error {
	log.Info("service delete")
	service := &core.Service{}
	err := json.Unmarshal([]byte(message), service)
	if err != nil {
		log.Errorf("unmarshal service error: %s", err.Error())
		return err
	}
	for _, port := range service.Spec.Ports {
		err = kubeproxy.DeleteService(service.Spec.ClusterIP, uint32(port.Port))
		if err != nil {
			log.Errorf("error in DeleteService err: %s", err.Error())
		}
	}
	FreeUsedIP(service.Spec.ClusterIP)
	err = DeleteEndpointObject(service)
	if err != nil {
		log.Errorf("error in DeleteEndpointObject")
		return err
	}
	return nil
}
