package service

import (
	"encoding/json"
	"errors"
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

	// creaete service and alloc ip
	clusterIP := FindUnusedIP()
	service.Spec.ClusterIP = clusterIP
	for _, port := range service.Spec.Ports {
		kubeproxy.CreateService(clusterIP, uint32(port.Port))
	}

	// get all pods
	response := utils.GetObject(core.ObjPod, service.MetaData.NameSpace, "")
	if response == "" {
		err = errors.New("cannot get pods")
		log.Errorf("get pod error: %s", err.Error())
		return err
	}
	pods := []core.Pod{}
	err = json.Unmarshal([]byte(response), &pods)
	if err != nil {
		log.Errorf("unmarshal pods error: %s", err.Error())
		return err
	}

	// select matched pods
	selectedPods := []core.Pod{}
	for _, pod := range pods {
		if MatchLabel(service.Spec.Selector, pod.MetaData.Labels) {
			selectedPods = append(selectedPods, pod)
		}
	}

	for _, port := range service.Spec.Ports {
		endpoint := kubeproxy.CreateEndpointObject(service, port.Port)
		err := utils.CreateObject(core.ObjEndPoint, endpoint.MetaData.Name, endpoint)
		if err != nil {
			log.Errorf("create endpoint error: %s", err.Error())
			return err
		}
		for _, pod := range selectedPods {
			destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
			kubeproxy.BindEndpoint(clusterIP, port.Port, pod.Status.PodIP, destPort)
			log.Infof("create endpoint: %s:%d -> %s:%d", clusterIP, port.Port, pod.Status.PodIP, destPort)
		}

	}
	return nil
}

func (sc *ServiceController) HandleUpdate(message string) error {
	//TODO implement me
	panic("implement me")
}

func (sc *ServiceController) HandleDelete(message string) error {
	//TODO implement me
	panic("implement me")
}
