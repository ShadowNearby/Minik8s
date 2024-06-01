package controller

import (
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"

	log "github.com/sirupsen/logrus"
)

type IController interface {
	GetChannel() string
	HandleCreate(message string) error
	HandleUpdate(message string) error
	HandleDelete(message string) error
	HandleTrigger(message string) error
}

func StartController(controller IController) {
	createChannel := constants.GenerateChannelName(controller.GetChannel(), constants.ChannelCreate)
	updateChannel := constants.GenerateChannelName(controller.GetChannel(), constants.ChannelUpdate)
	deleteChannel := constants.GenerateChannelName(controller.GetChannel(), constants.ChannelDelete)
	triggerChannel := constants.GenerateChannelName(controller.GetChannel(), constants.ChannelTrigger)
	createMessages := storage.RedisInstance.SubscribeChannel(createChannel)
	updateMessages := storage.RedisInstance.SubscribeChannel(updateChannel)
	deleteMessages := storage.RedisInstance.SubscribeChannel(deleteChannel)
	triggerMessages := storage.RedisInstance.SubscribeChannel(triggerChannel)
	go func() {
		for {
			for message := range createMessages {
				err := controller.HandleCreate(message.Payload)
				if err != nil {
					log.Errorf("handle create error: %s", err.Error())
				}
			}
		}
	}()
	go func() {
		for {
			for message := range updateMessages {
				err := controller.HandleUpdate(message.Payload)
				if err != nil {
					log.Errorf("handle update error: %s", err.Error())
				}
			}
		}
	}()
	go func() {
		for {
			for message := range deleteMessages {
				err := controller.HandleDelete(message.Payload)
				if err != nil {
					log.Errorf("handle delete error: %s", err.Error())
				}
			}
		}
	}()
	go func() {
		for {
			for message := range triggerMessages {
				err := controller.HandleTrigger(message.Payload)
				if err != nil {
					log.Errorf("handle delete error: %s", err.Error())
				}
			}
		}
	}()

}
