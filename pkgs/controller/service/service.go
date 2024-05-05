package service

import (
	"encoding/json"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/kubeproxy"

	log "github.com/sirupsen/logrus"
)

const TotalIP = (1 << 8)

var UsedIP = [TotalIP]bool{}
var ServiceSelector = map[string]*core.Selector{}

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

	// creaete service and alloc ip
	clusterIP := FindUnusedIP()
	service.Spec.ClusterIP = clusterIP
	controller.SetObject(core.ObjService, service.MetaData.NameSpace, service.MetaData.Name, service)
	for _, port := range service.Spec.Ports {
		kubeproxy.CreateService(clusterIP, uint32(port.Port))
	}
	PutSelector(service)

	err = CreateEndpointObject(service)
	if err != nil {
		log.Errorf("error in CreateEndpointObject")
		return err
	}
	return nil
}

func (sc *ServiceController) HandleUpdate(message string) error {
	service := &core.Service{}
	err := json.Unmarshal([]byte(message), service)
	if err != nil {
		log.Errorf("unmarshal service error: %s", err.Error())
		return err
	}
	previousSelector := GetSelector(service)
	if MatchLabel(previousSelector.MatchLabels, service.Spec.Selector.MatchLabels) {
		return nil
	}

	err = DeleteEndpointObject(service, nil)
	if err != nil {
		log.Errorf("error in UpdateEndpointObject")
		return err
	}
	err = CreateEndpointObject(service)
	if err != nil {
		log.Errorf("error in UpdateEndpointObject")
		return err
	}

	PutSelector(service)
	return nil
}

func (sc *ServiceController) HandleDelete(message string) error {
	service := &core.Service{}
	err := json.Unmarshal([]byte(message), service)
	if err != nil {
		log.Errorf("unmarshal service error: %s", err.Error())
		return err
	}

	DelSelector(service)

	FreeUsedIP(service.Spec.ClusterIP)
	err = DeleteEndpointObject(service, nil)
	if err != nil {
		log.Errorf("error in DeleteEndpointObject")
		return err
	}
	return nil
}
