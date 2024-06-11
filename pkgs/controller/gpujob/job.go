package gpujob

import (
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/controller/tools"
	"minik8s/utils"
)

type JobController struct {
}

func (r *JobController) GetChannel() string {
	return constants.ChannelJob
}

func (r *JobController) HandleCreate(message string) error {
	var job core.Job
	_ = utils.JsonUnMarshal(message, &job)
	tools.CreatePodforJob(job)

	res := utils.GetObject(core.ObjPod, job.GetNamespace(), job.MetaData.Name)
	var pod core.Pod
	_ = utils.JsonUnMarshal(res, &pod)
	job.Status = pod.Status
	job.Status.Phase = core.PodPhaseRunning
	_ = utils.SetObject(core.ObjJob, job.GetNamespace(), job.MetaData.Name, job)

	log.Info("[job controller] Create job. Name: ", job.MetaData.Name)
	return nil

}

func (r *JobController) HandleDelete(message string) error {
	var job core.Job
	_ = utils.JsonUnMarshal(message, job)
	tools.DeletePodforJob(job)
	_ = utils.DeleteObject(core.ObjJob, job.MetaData.Namespace, job.MetaData.Name)
	log.Info("[job controller] Delete job. Name:", job.MetaData.Name)
	return nil
}

func (r *JobController) HandleUpdate(message string) error {
	return nil
}
