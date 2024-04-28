package rsc

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/config"
	"minik8s/pkgs/controller"
	"minik8s/utils"
	"sync"
)

type ReplicaSetController struct {
	// subscribed channels
	UpdateReplicaChan <-chan *redis.Message
	NewReplicaChan    <-chan *redis.Message
}

var ReplicaSetControllerInstance ReplicaSetController

func (rsc *ReplicaSetController) Run() error {
	rsc.UpdateReplicaChan = storage.RedisInstance.SubscribeChannel(storage.ChannelUpdateReplica)
	rsc.NewReplicaChan = storage.RedisInstance.SubscribeChannel(storage.ChannelNewReplica)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(g *sync.WaitGroup) {
		defer g.Done()
		rsc.listenNewRS()
	}(&wg)
	go func(g *sync.WaitGroup) {
		defer g.Done()
		rsc.listenUpdateRS()
	}(&wg)
	wg.Wait()
	return errors.New("replicaset controller unexpect exit")
}

func (rsc *ReplicaSetController) listenNewRS() {
	for {
		for msg := range rsc.NewReplicaChan {
			fmt.Println(msg.Channel, msg.Payload)
			err := rsc.syncReplicas(msg.Payload)
			if err != nil {
				logger.Errorf("sync new replica error: %s", err.Error())
			}
		}
	}
}

func (rsc *ReplicaSetController) listenUpdateRS() {
	for {
		for msg := range rsc.UpdateReplicaChan {
			fmt.Println(msg.Channel, msg.Payload)
			err := rsc.syncReplicas(msg.Payload)
			if err != nil {
				logger.Errorf("sync update replica error: %s", err.Error())
			}
		}
	}
}

func (rsc *ReplicaSetController) syncReplicas(key string) error {
	namespace, name, err := controller.SplitChannelInfo(key)
	if err != nil {
		return err
	}
	replicaTxt := controller.GetObject(config.ObjReplicaSet, namespace, name)
	if replicaTxt == "" {
		return errors.New("cannot get replica")
	}
	var replica core.ReplicaSet
	err = utils.JsonUnMarshal(replicaTxt, &replica)
	if err != nil {
		return err
	}
	err = rsc.manageReplicas(&replica)
	if err != nil {
		return err
	}
	return nil
}

func (rsc *ReplicaSetController) manageReplicas(rs *core.ReplicaSet) error {
	// first get pods within the rsc namespace
	var pods = make([]core.Pod, 0)
	podsListTxt := controller.GetObject(config.ObjPod, rs.MetaData.NameSpace, "")
	if podsListTxt == "" {
		return errors.New("cannot get pods")
	}
	err := utils.JsonUnMarshal(podsListTxt, &pods)
	if err != nil {
		return err
	}
	// second filter the pods meets selector and don't have controller
	targets := rsc.filterPods(&pods, rs)
	// if not enough, create new pods
	if len(targets) < rs.Spec.Replicas {
		pod := core.Pod{
			ApiVersion: rs.ApiVersion,
			MetaData:   rs.Spec.Template.MetaData,
			Spec:       rs.Spec.Template.Spec,
			Status:     core.PodStatus{},
		}
		setController(&pod, rs)
		ops := rs.Spec.Replicas - len(targets)
		for i := 0; i < ops; i++ {
			err = controller.CreateObject(config.ObjPod, rs.MetaData.NameSpace, pod)
			if err != nil {
				return err
			}
			targets = append(targets, pod)
		}
	}
	return nil
}

func (rsc *ReplicaSetController) filterPods(origin *[]core.Pod, rs *core.ReplicaSet) []core.Pod {
	result := make([]core.Pod, 0)
	for _, pod := range *origin {
		if len(result) >= rs.Spec.Replicas {
			break
		}
		flag := true
		if pod.MetaData.OwnerReference.Controller == false {
			for key, value := range rs.MetaData.Labels {
				val, ok := pod.MetaData.Labels[key]
				if !ok || val != value {
					flag = false
					break
				}
			}
			if flag == true {
				setController(&pod, rs)
				result = append(result, pod)
			}
		}
	}
	return result
}

func setController(pod *core.Pod, rs *core.ReplicaSet) {
	pod.MetaData.OwnerReference.Controller = true
	pod.MetaData.OwnerReference.ObjType = config.ObjReplicaSet
	pod.MetaData.OwnerReference.Name = rs.MetaData.Name
}
