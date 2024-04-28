package controller

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
)

var UsedIP = map[string]bool{}

type IServiceController struct {
}

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
	if UsedIP[service.Spec.ClusterIP] == true {
		log.Errorf("ip %s is already used", service.Spec.ClusterIP)
		return errors.New("ip is already used")
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
