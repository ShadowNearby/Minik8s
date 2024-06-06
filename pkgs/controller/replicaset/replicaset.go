package rsc

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"time"

	logger "github.com/sirupsen/logrus"
)

type ReplicaSetController struct {
}

func (rsc *ReplicaSetController) GetChannel() string {
	return constants.ChannelReplica
}

func (rsc *ReplicaSetController) BackGroundTask() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		var replicas []core.ReplicaSet
		rsTxt := utils.GetObject(core.ObjReplicaSet, "default", "")
		utils.JsonUnMarshal(rsTxt, &replicas)
		for _, replica := range replicas {
			pods, err := utils.FindRSPods(replica.MetaData.Name)
			if err != nil {
				continue
			}
			replica.Status.RealReplicas = len(pods)
			if replica.Status.RealReplicas > replica.Spec.Replicas {
				rsc.scaleDown(&replica, pods)
			} else if replica.Status.RealReplicas < replica.Spec.Replicas {
				rsc.scaleUp(&replica, pods)
			}
		}
	}
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
	// if the rs has owner, we don't need to manage it
	if replica.MetaData.OwnerReference.Controller {
		return nil
	}
	pods, _ := utils.FindRSPods(replica.MetaData.Name)
	replica.Spec.Replicas = 0
	err = rsc.scaleDown(&replica, pods)
	return err
	//err = storage.Del(fmt.Sprintf("/replicas/object/%s/%s", replica.MetaData.NameSpace, replica.MetaData.Name))
}

func (rsc *ReplicaSetController) createReplicas(info string) error {
	var replica core.ReplicaSet
	err := utils.JsonUnMarshal(info, &replica)
	if err != nil {
		return err
	}
	target := make([]core.Pod, 0)
	err = rsc.scaleUp(&replica, target)
	if err != nil {
		return err
	}
	// update the replicaset status
	err = utils.SetObjectStatus(core.ObjReplicaSet, replica.MetaData.Namespace, replica.MetaData.Name, replica)
	return err
}

func (rsc *ReplicaSetController) updateReplicas(info string) error {
	var replicas []core.ReplicaSet
	err := utils.JsonUnMarshal(info, &replicas)

	if err != nil {
		logger.Errorf("unmarshal replicas error: %s", err.Error())
		return err
	}
	if len(replicas) < 2 {
		return fmt.Errorf("not enough replica info")
	}
	err = rsc.manageUpdateReplicas(&replicas[0], &replicas[1])
	if err != nil {
		return err
	}
	err = utils.SetObjectStatus(core.ObjReplicaSet, "default", replicas[1].MetaData.Name, replicas[1])
	return err
}

func (rsc *ReplicaSetController) manageUpdateReplicas(oldRs *core.ReplicaSet, newRs *core.ReplicaSet) error {
	// firstly, new rs should inherent old rs state
	newRs.Status = oldRs.Status
	// if target replica changes we should add or delete pods
	// if template changes we should delete all existed pods and create new (we don't consider this case currently)

	// get pods managed by this replica
	pods, err := utils.FindRSPods(newRs.MetaData.Name)
	if err != nil {
		logger.Errorf("failed getting rs pods")
		return err
	}
	newRs.Status.RealReplicas = len(pods)
	if len(pods) > newRs.Spec.Replicas {
		// delete pods
		for _, pod := range pods[newRs.Spec.Replicas:] {
			err := utils.DeleteObject(core.ObjPod, pod.MetaData.Namespace, pod.MetaData.Name)
			if err != nil {
				logger.Errorf("delete pod error: %s", err.Error())
			}
			newRs.Status.RealReplicas--
		}
	} else if len(pods) < newRs.Spec.Replicas {
		// create pods
		pod := generateRSPod(newRs)
		templateName := pod.MetaData.Name
		setController(&pod, newRs)
		ops := newRs.Spec.Replicas - len(pods)
		for i := 0; i < ops; i++ {
			// should regenerate pod uuid
			pod.MetaData.Name = fmt.Sprintf("rs-%s-%s", templateName, utils.GenerateUUID(6))
			logger.Infof("pod name: %s", pod.MetaData.Name)
			pod.MetaData.UUID = utils.GenerateUUID()
			err = utils.CreateObject(core.ObjPod, newRs.MetaData.Namespace, pod)
			if err != nil {
				return err
			}
			newRs.Status.RealReplicas++
		}
	}
	logger.Infof("updated replicas, real replica: %d, spec replica: %d", newRs.Status.RealReplicas, newRs.Spec.Replicas)

	return nil
}

func (rsc *ReplicaSetController) scaleDown(rs *core.ReplicaSet, targets []core.Pod) error {
	ops := len(targets) - rs.Spec.Replicas
	rs.Status.RealReplicas = len(targets)
	for i := 0; i < ops; i++ {
		target := targets[i]
		err := utils.DeleteObject(core.ObjPod, target.MetaData.Namespace, target.MetaData.Name)
		if err != nil {
			logger.Errorf("delete object error: %s", err.Error())
			return err
		}
		rs.Status.RealReplicas--
	}
	return nil
}

func (rsc *ReplicaSetController) scaleUp(rs *core.ReplicaSet, targets []core.Pod) error {
	// targets := make([]core.Pod, 0)
	if len(targets) < rs.Spec.Replicas {
		pod := generateRSPod(rs)
		templateName := pod.MetaData.Name
		setController(&pod, rs)
		ops := rs.Spec.Replicas - len(targets)
		for i := 0; i < ops; i++ {
			// should re-generate pod uuid
			pod.MetaData.Name = fmt.Sprintf("rs-%s-%s", templateName, utils.GenerateUUID(6))
			logger.Infof("pod name: %s", pod.MetaData.Name)
			pod.MetaData.UUID = utils.GenerateUUID()
			err := utils.CreateObject(core.ObjPod, rs.MetaData.Namespace, pod)
			if err != nil {
				return err
			}
			rs.Status.RealReplicas++
		}
	}
	return nil
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
}

func generateRSPod(rs *core.ReplicaSet) core.Pod {
	pod := core.Pod{
		ApiVersion: rs.ApiVersion,
		MetaData:   rs.Spec.Template.MetaData,
		Spec:       rs.Spec.Template.Spec,
		Status:     core.PodStatus{},
	}
	// set random replica name and namespace
	//pod.MetaData.Name = fmt.Sprintf("rs-%s-%s", pod.MetaData.Name, utils.GenerateUUID(6))
	pod.MetaData.Namespace = "default"
	pod.MetaData.UUID = utils.GenerateUUID()
	setController(&pod, rs)
	return pod
}
