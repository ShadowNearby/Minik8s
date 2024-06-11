package tools

import (
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

func CreatePodforJob(job core.Job) {
	pod := &core.Pod{
		ApiVersion: job.ApiVersion,
		MetaData:   job.MetaData,
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
	log.Info("[job controller] create pod for job: \n", pod)
	_ = utils.CreateObject(core.ObjPod, job.MetaData.Namespace, pod)
}
func DeletePodforJob(job core.Job) {
	_ = utils.DeleteObject(core.ObjPod, job.MetaData.Namespace, job.MetaData.Name)
}
