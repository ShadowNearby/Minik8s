package podcontroller

import (
	"encoding/json"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/controller/tools"
	"minik8s/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

type PodController struct{}

func (pc *PodController) GetChannel() string {
	return constants.ChannelPod
}

func (pc *PodController) HandleCreate(message string) error {
	return nil
}

func (pc *PodController) HandleUpdate(message string) error {
	var pod core.Pod
	_ = utils.JsonUnMarshal(message, pod)
	if pod.MetaData.OwnerReference.ObjType == core.ObjJob {
		log.Info("[job controller] handle pdate")
		info := utils.GetObject(core.ObjJob, pod.MetaData.Namespace, pod.MetaData.OwnerReference.Name)
		var job core.Job
		_ = utils.JsonUnMarshal(info, job)

		job.Status.Phase = pod.Status.Phase
		switch pod.Status.Phase {
		case core.PodPhaseFailed:
			{
				job.Spec.BackoffLimit -= 1
				if job.Spec.BackoffLimit > 0 {
					go func(job core.Job) {
						tools.DeletePodforJob(job)
						time.Sleep(time.Second * 5)
						tools.CreatePodforJob(job)
					}(job)
				}
			}
		case core.PodPhaseSucceeded:
			{
				waitToDelete := func(t int, job core.Job) {
					time.Sleep(time.Second * time.Duration(t))
					tools.DeletePodforJob(job)
				}
				go waitToDelete(job.Spec.TtlSecondsAfterFinished, job)
			}
		}
		log.Info("[job controller] update phase:", job.Status.Phase)
		_ = utils.SetObject(core.ObjJob, job.MetaData.Namespace, pod.MetaData.Name, job)
	}
	return nil
}

func (pc *PodController) HandleDelete(message string) error {
	pod := &core.Pod{}
	err := json.Unmarshal([]byte(message), pod)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}
	log.Infof("delete pod: %s:%s", pod.MetaData.Namespace, pod.MetaData.Name)
	_, _, err = utils.SendRequest("DELETE", fmt.Sprintf("http://%s:%s/pod/stop/%s/%s", pod.Status.HostIP, config.NodePort, pod.GetNamespace(), pod.MetaData.Name), nil)
	if err != nil {
		log.Errorf("delete pod error: %s", err.Error())
		return err
	}
	return nil
}
