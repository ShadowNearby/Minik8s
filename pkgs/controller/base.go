package controller

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
)

type IController interface {
	GetChannel() string
	HandleCreate(message string) error
	HandleUpdate(message string) error
	HandleDelete(message string) error
}

func StartController(controller IController) {
	createChannel := fmt.Sprintf("%s-%s", controller.GetChannel(), constants.ChannelCreate)
	updateChannel := fmt.Sprintf("%s-%s", controller.GetChannel(), constants.ChannelUpdate)
	deleteChannel := fmt.Sprintf("%s-%s", controller.GetChannel(), constants.ChannelDelete)
	createMessages := storage.RedisInstance.SubscribeChannel(createChannel)
	updateMessages := storage.RedisInstance.SubscribeChannel(updateChannel)
	deleteMessages := storage.RedisInstance.SubscribeChannel(deleteChannel)
	go func() {
		for message := range createMessages {
			err := controller.HandleCreate(message.Payload)
			if err != nil {
				log.Errorf("handle create error: %s", err.Error())
			}
		}
	}()
	go func() {
		for message := range updateMessages {
			err := controller.HandleCreate(message.Payload)
			if err != nil {
				log.Errorf("handle update error: %s", err.Error())
			}
		}
	}()
	go func() {
		for message := range deleteMessages {
			err := controller.HandleCreate(message.Payload)
			if err != nil {
				log.Errorf("handle delete error: %s", err.Error())
			}
		}
	}()
}
