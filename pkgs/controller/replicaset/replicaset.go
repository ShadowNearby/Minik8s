package rsc

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"
)

type ReplicaSetController struct {
}

func (rsc *ReplicaSetController) GetChannel() string {
	return constants.ChannelReplica
}

func (rsc *ReplicaSetController) HandleCreate(key string) error {
	err := rsc.createReplicas(key)
	if err != nil {
		return err
	}
	return nil
}
func (rsc *ReplicaSetController) HandleUpdate(key string) error {
	err := rsc.updateReplicas(key)
	return err
}
func (rsc *ReplicaSetController) HandleDelete(key string) error {
	err := rsc.deleteReplicas(key)
	return err
}

func (rsc *ReplicaSetController) deleteReplicas(info string) error {
	var replica core.ReplicaSet
	err := utils.JsonUnMarshal(info, &replica)
	if err != nil {
		return err
	}
	err = rsc.manageDelReplicas(&replica)
	return err
	//err = storage.Del(fmt.Sprintf("/replicas/object/%s/%s", replica.MetaData.NameSpace, replica.MetaData.Name))
}

func (rsc *ReplicaSetController) createReplicas(info string) error {
	var replica core.ReplicaSet
	err := utils.JsonUnMarshal(info, &replica)
	if err != nil {
		return err
	}
	err = rsc.manageCreateReplicas(&replica)
	if err != nil {
		return err
	}
	// update the replicaset status
	err = utils.SetObject(core.ObjReplicaSet, replica.MetaData.NameSpace, replica.MetaData.Name, replica)
	return err
}

func (rsc *ReplicaSetController) updateReplicas(info string) error {
	var replicas []core.ReplicaSet
	err := utils.JsonUnMarshal(info, &replicas)

	if err != nil {
		return err
	}
	if len(replicas) < 2 {
		return fmt.Errorf("not enough replica info")
	}
	err = rsc.manageUpdateReplicas(&replicas[0], &replicas[1])
	if err != nil {
		return err
	}
	err = utils.SetObject(core.ObjReplicaSet, replicas[1].MetaData.NameSpace, replicas[1].MetaData.Name, replicas[1])
	return err
}

func (rsc *ReplicaSetController) manageUpdateReplicas(oldRs *core.ReplicaSet, newRs *core.ReplicaSet) error {
	if oldRs.MetaData.NameSpace != newRs.MetaData.NameSpace {
		err := rsc.manageDelReplicas(oldRs)
		if err != nil {
			return err
		}
		err = rsc.manageCreateReplicas(newRs)
		if err != nil {
			return err
		}
	} else {
		var pods = make([]core.Pod, 0)
		podsListTxt := utils.GetObject(core.ObjPod, newRs.MetaData.NameSpace, "")
		if podsListTxt == "" {
			return errors.New("cannot get pods")
		}
		err := utils.JsonUnMarshal(podsListTxt, &pods)
		if err != nil {
			return err
		}
		targets := rsc.filterOwners(&pods, newRs)
		newRs.Status.RealReplicas = len(targets)
		if len(targets) > newRs.Spec.Replicas {
			// delete pods
			for _, pod := range targets[newRs.Spec.Replicas:] {
				err := utils.DeleteObject(core.ObjPod, pod.MetaData.NameSpace, pod.MetaData.Name)
				if err != nil {
					logger.Errorf("delete pod error: %s", err.Error())
				}
				newRs.Status.RealReplicas--
			}
		} else if len(targets) < newRs.Spec.Replicas {
			// create pods
			pod := core.Pod{
				ApiVersion: newRs.ApiVersion,
				MetaData:   newRs.MetaData,
				Spec:       newRs.Spec.Template.Spec,
				Status:     core.PodStatus{},
			}
			setController(&pod, newRs)
			ops := newRs.Spec.Replicas - len(targets)
			for i := 0; i < ops; i++ {
				err = utils.CreateObject(core.ObjPod, newRs.MetaData.NameSpace, pod)
				if err != nil {
					return err
				}
				newRs.Status.RealReplicas++
			}
		}
	}
	logger.Infof("updated replicas, real replica: %d, spec replica: %d", newRs.Status.RealReplicas, newRs.Spec.Replicas)

	return nil
}

func (rsc *ReplicaSetController) manageDelReplicas(rs *core.ReplicaSet) error {
	var pods = make([]core.Pod, 0)
	podListTxt := utils.GetObject(core.ObjPod, rs.MetaData.NameSpace, "")
	if podListTxt == "" {
		return errors.New("cannot get pods")
	}
	err := utils.JsonUnMarshal(podListTxt, &pods)
	if err != nil {
		return err
	}
	targets := rsc.filterOwners(&pods, rs)
	rs.Status.RealReplicas = len(targets)
	for _, target := range targets {
		err = utils.DeleteObject(core.ObjPod, target.MetaData.NameSpace, target.MetaData.Name)
		if err != nil {
			logger.Errorf("delete object error: %s", err.Error())
			return err
		}
		rs.Status.RealReplicas--
	}
	return nil
}

func (rsc *ReplicaSetController) manageCreateReplicas(rs *core.ReplicaSet) error {
	// first get pods within the rsc namespace
	var pods = make([]core.Pod, 0)
	podsListTxt := utils.GetObject(core.ObjPod, rs.MetaData.NameSpace, "")
	if podsListTxt == "" {
		return errors.New("cannot get pods")
	}
	err := utils.JsonUnMarshal(podsListTxt, &pods)
	if err != nil {
		return err
	}
	// second filter the pods meets selector and don't have controller
	targets := rsc.selectPods(&pods, rs)
	rs.Status.RealReplicas = len(targets)
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
			err = utils.CreateObject(core.ObjPod, rs.MetaData.NameSpace, pod)
			if err != nil {
				return err
			}
			rs.Status.RealReplicas++
		}
	}
	return nil
}

func (rsc *ReplicaSetController) filterOwners(origin *[]core.Pod, rs *core.ReplicaSet) []core.Pod {
	result := make([]core.Pod, 0)
	for _, pod := range *origin {
		or := pod.MetaData.OwnerReference
		if or.Controller == true &&
			or.ObjType == core.ObjReplicaSet &&
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
	pod.MetaData.OwnerReference.ObjType = core.ObjReplicaSet
	pod.MetaData.OwnerReference.Name = rs.MetaData.Name
	pod.MetaData.OwnerReference.NameSpace = rs.MetaData.NameSpace
}
