package kubeproxy

import (
	"encoding/json"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type ServiceController struct{}

func (sc *ServiceController) Run() {
	var redisInstance = &storage.Redis{
		Client:   storage.CreateRedisClient(),
		Channels: make(map[string]*redis.PubSub),
	}

	createChannel := constants.GenerateChannelName(sc.GetChannel(), constants.ChannelCreate)
	updateChannel := constants.GenerateChannelName(sc.GetChannel(), constants.ChannelUpdate)
	deleteChannel := constants.GenerateChannelName(sc.GetChannel(), constants.ChannelDelete)

	redisInstance.CreateChannel(createChannel)
	redisInstance.CreateChannel(updateChannel)
	redisInstance.CreateChannel(deleteChannel)

	createMessages := redisInstance.SubscribeChannel(createChannel)
	updateMessages := redisInstance.SubscribeChannel(updateChannel)
	deleteMessages := redisInstance.SubscribeChannel(deleteChannel)

	go func() {
		for {
			for message := range createMessages {
				err := sc.HandleCreate(message.Payload)
				if err != nil {
					log.Errorf("handle create error: %s", err.Error())
				}
			}
		}
	}()
	go func() {
		for {
			for message := range updateMessages {
				err := sc.HandleUpdate(message.Payload)
				if err != nil {
					log.Errorf("handle update error: %s", err.Error())
				}
			}
		}
	}()
	go func() {
		for {
			for message := range deleteMessages {
				err := sc.HandleDelete(message.Payload)
				if err != nil {
					log.Errorf("handle delete error: %s", err.Error())
				}
			}
		}
	}()
}

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
	if service.Spec.ClusterIP == "" && service.Spec.Type == core.ServiceTypeClusterIP {
		clusterIP = FindUnusedIP(service.MetaData.Namespace, service.MetaData.Name)
		service.Spec.ClusterIP = clusterIP
		utils.SetObject(core.ObjService, service.MetaData.Namespace, service.MetaData.Name, service)
	}

	// creaete service and alloc ip
	if service.Spec.Type == core.ServiceTypeClusterIP {
		for _, port := range service.Spec.Ports {
			CreateService(clusterIP, port.Port)
		}
	} else if service.Spec.Type == core.ServiceTypeNodePort {
		for _, port := range service.Spec.Ports {
			CreateService(constants.AllIP, port.NodePort)
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
	service := &core.Service{}
	err := json.Unmarshal([]byte(message), service)
	if err != nil {
		log.Errorf("unmarshal service error: %s", err.Error())
		return err
	}
	if service.Spec.Type == core.ServiceTypeClusterIP {
		for _, port := range service.Spec.Ports {
			err = DeleteService(service.Spec.ClusterIP, uint32(port.Port))
			if err != nil {
				log.Errorf("error in DeleteService err: %s", err.Error())
			}
		}
		FreeUsedIP(service.MetaData.Namespace, service.MetaData.Name)
	} else if service.Spec.Type == core.ServiceTypeNodePort {
		for _, port := range service.Spec.Ports {
			err = DeleteService(constants.AllIP, uint32(port.NodePort))
			if err != nil {
				log.Errorf("error in DeleteService err: %s", err.Error())
			}
		}
	}
	err = DeleteEndpointObject(service)
	if err != nil {
		log.Errorf("error in DeleteEndpointObject")
		return err
	}
	return nil
}
