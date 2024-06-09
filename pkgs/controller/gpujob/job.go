package gpujob

import (
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"time"
)

type JobController struct {
}

type jobPodHandler struct {
}

func (r *JobController) GetChannel() string {
	return constants.ChannelJob
}

func (r *JobController) HandleCreate(message string) error {
	var job core.Job
	_ = utils.JsonUnMarshal(message, &job)
	createPodforJob(job)
	job.Status.Phase = core.PodPhaseRunning
	_ = utils.SetObject(core.ObjJob, job.GetNamespace(), job.MetaData.Name, job)
	log.Info("[job controller] Create job. Name: ", job.MetaData.Name)
	return nil
}

func (r *JobController) HandleDelete(message string) error {
	var job core.Job
	_ = utils.JsonUnMarshal(message, job)
	deletePodforJob(job)
	_ = utils.DeleteObject(core.ObjJob, job.MetaData.Namespace, job.MetaData.Name)
	log.Info("[job controller] Delete job. Name:", job.MetaData.Name)
	return nil
}

func (r *JobController) HandleUpdate(message string) error {
	return nil
}

/* ========== Start Pod Handler ========== */

func (p jobPodHandler) HandleUpdate(message string) {
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
						deletePodforJob(job)
						time.Sleep(time.Second * 5)
						createPodforJob(job)
					}(job)
				}
			}
		case core.PodPhaseSucceeded:
			{
				waitToDelete := func(t int, job core.Job) {
					time.Sleep(time.Second * time.Duration(t))
					deletePodforJob(job)
				}
				go waitToDelete(job.Spec.TtlSecondsAfterFinished, job)
			}
		}
		log.Info("[job controller] update phase:", job.Status.Phase)
		_ = utils.SetObject(core.ObjJob, job.MetaData.Namespace, pod.MetaData.Name, job)
	}
}

func createPodforJob(job core.Job) {
	pod := &core.Pod{
		MetaData: job.MetaData,
		Spec: core.PodSpec{
			Selector: core.Selector{
				MatchLabels: job.Spec.NodeSelector,
			},
			Containers: job.Spec.Containers,
			Volumes:    job.Spec.Volumes,
		},
	}
	pod.MetaData.OwnerReference = core.OwnerReference{
		ObjType:    core.ObjJob,
		Name:       job.MetaData.Name,
		Controller: false,
	}
	_ = utils.CreateObject(core.ObjPod, job.MetaData.Namespace, pod)
}
func deletePodforJob(job core.Job) {
	_ = utils.DeleteObject(core.ObjPod, job.MetaData.Namespace, job.MetaData.Name)
}
