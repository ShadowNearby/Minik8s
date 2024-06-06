package kubeproxy

import (
	"encoding/json"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type EndpointController struct {
}

func (sc *EndpointController) Run() {
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

func (sc *EndpointController) GetChannel() string {
	// Pod 修改时需要修改Endpoint
	return constants.ChannelEndpoint
}

func (sc *EndpointController) HandleCreate(message string) error {
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

	if prePod.Status.PodIP == pod.Status.PodIP && prePod.Status.Phase == core.PodPhaseRunning && pod.Status.Phase == core.PodPhaseRunning {
		log.Info("podIP not changed & pod is still running")
		return nil
	}

	services, err := GetAllServiceObject(pod.MetaData.Namespace)
	if err != nil {
		log.Errorf("get all service error: %s", err.Error())
		return err
	}
	for _, service := range services {
		if prePod.Status.PodIP != pod.Status.PodIP {
			if prePod.Status.PodIP != "" {
				UpdateEndpointObjectByPodDelete(&service, &prePod)
			}
			UpdateEndpointObjectByPodCreate(&service, &pod)
		}
		if prePod.Status.Phase == core.PodPhaseRunning && pod.Status.Phase != core.PodPhaseRunning {
			UpdateEndpointObjectByPodDelete(&service, &pod)
		}
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
