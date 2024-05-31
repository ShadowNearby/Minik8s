package service

import (
	"encoding/json"
	"errors"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"

	log "github.com/sirupsen/logrus"
)

type EndpointController struct {
}

func (sc *EndpointController) GetChannel() string {
	// Pod 修改时需要修改Endpoint
	return constants.ChannelPod
}
func (sc *EndpointController) HandleTrigger(message string) error {
	log.Errorf("EndpointController cannot be triggered")
	return errors.New("EndpointController cannot be triggered")
}
func (sc *EndpointController) HandleCreate(message string) error {
	pod := &core.Pod{}
	err := json.Unmarshal([]byte(message), pod)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}

	services, err := GetAllServiceObject(pod.MetaData.Namespace)
	if err != nil {
		log.Errorf("get all service error: %s", err.Error())
		return err
	}
	for _, service := range services {
		UpdateEndpointObjectByPodCreate(&service, pod)
	}
	return nil
}

func (sc *EndpointController) HandleUpdate(message string) error {
	pods := []core.Pod{}
	err := json.Unmarshal([]byte(message), &pods)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}
	if len(pods) != 2 {
		return fmt.Errorf("endpoint update error")
	}
	prePod := pods[0]
	pod := pods[1]

	services, err := GetAllServiceObject(pod.MetaData.Namespace)
	if err != nil {
		log.Errorf("get all service error: %s", err.Error())
		return err
	}
	for _, service := range services {
		UpdateEndpointObjectByPodDelete(&service, &prePod)
		UpdateEndpointObjectByPodCreate(&service, &pod)
	}
	return nil
}

func (sc *EndpointController) HandleDelete(message string) error {
	pod := &core.Pod{}
	err := json.Unmarshal([]byte(message), pod)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}

	services, err := GetAllServiceObject(pod.MetaData.Namespace)
	if err != nil {
		log.Errorf("get all service error: %s", err.Error())
		return err
	}
	for _, service := range services {
		UpdateEndpointObjectByPodDelete(&service, pod)
	}
	return nil
}
