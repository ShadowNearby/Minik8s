package service

import (
	"encoding/json"
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

func (sc *EndpointController) HandleCreate(message string) error {
	return nil
}

func (sc *EndpointController) HandleUpdate(message string) error {
	pod := &core.Service{}
	err := json.Unmarshal([]byte(message), pod)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}

	// services, err := GetAllServiceObject(pod.MetaData.NameSpace)
	// for _, service := range services {

	// }
	return nil
}

func (sc *EndpointController) HandleDelete(message string) error {
	pod := &core.Service{}
	err := json.Unmarshal([]byte(message), pod)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}
	return nil
}
