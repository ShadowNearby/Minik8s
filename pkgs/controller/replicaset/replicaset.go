package rsc

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/config"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/controller"
	"minik8s/utils"
)

type ReplicaSetController struct {
}

func (rsc *ReplicaSetController) GetChannel() string {
	return constants.ChannelReplica
}

func (rsc *ReplicaSetController) HandleCreate(key string) error {
	err := rsc.syncReplicas(key)
	if err != nil {
		return err
	}
	return nil
}
func (rsc *ReplicaSetController) HandleUpdate(key string) error {
	err := rsc.syncReplicas(key)
	return err
}
func (rsc *ReplicaSetController) HandleDelete(key string) error {
	err := rsc.deleteReplicas(key)
	return err
}

func (rsc *ReplicaSetController) deleteReplicas(key string) error {
	namespace, name, err := controller.SplitChannelInfo(key)
	if err != nil {
		return err
	}
	replicaTxt := controller.GetObject(config.ObjReplicaSet, namespace, name)
	if replicaTxt == "" {
		return errors.New("cannot get replica")
	} /* delete in storage after clear pods*/
	var replica core.ReplicaSet
	err = utils.JsonUnMarshal(replicaTxt, &replica)
	if err != nil {
		return err
	}
	err = storage.Del(fmt.Sprintf("/replicas/object/%s/%s", namespace, name))
	return err
}

func (rsc *ReplicaSetController) syncReplicas(key string) error {
	namespace, name, err := controller.SplitChannelInfo(key)
	if err != nil {
		return err
	}
	replicaTxt := controller.GetObject(config.ObjReplicaSet, namespace, name)
	if replicaTxt == "" {
		return errors.New("cannot get replica")
	} /* redis and etcd has store the replica */
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

func (rsc *ReplicaSetController) manageDelReplicas(rs *core.ReplicaSet) error {
	var pods = make([]core.Pod, 0)
	podListTxt := controller.GetObject(config.ObjPod, rs.MetaData.NameSpace, "")
	if podListTxt == "" {
		return errors.New("cannot get pods")
	}
	err := utils.JsonUnMarshal(podListTxt, &pods)
	if err != nil {
		return err
	}
	targets := rsc.filterOwners(&pods, rs)
	for _, target := range targets {
		err = controller.DeleteObject(config.ObjPod, target.MetaData.NameSpace, target.MetaData.Name)
		if err != nil {
			logger.Errorf("delete object error: %s", err.Error())
			return err
		}
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
	targets := rsc.selectPods(&pods, rs)
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

func (rsc *ReplicaSetController) filterOwners(origin *[]core.Pod, rs *core.ReplicaSet) []core.Pod {
	result := make([]core.Pod, 0)
	for _, pod := range *origin {
		or := pod.MetaData.OwnerReference
		if or.Controller == true &&
			or.ObjType == config.ObjReplicaSet &&
			or.NameSpace == rs.MetaData.NameSpace &&
			or.Name == rs.MetaData.Name {
			result = append(result, pod)
		}
	}
	return result
}

func (rsc *ReplicaSetController) selectPods(origin *[]core.Pod, rs *core.ReplicaSet) []core.Pod {
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
	pod.MetaData.OwnerReference.NameSpace = rs.MetaData.NameSpace
}
