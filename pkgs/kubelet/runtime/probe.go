package runtime

import (
	"errors"
	"github.com/containerd/containerd"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"
)

type ProbeType string

const (
	ExecProbe ProbeType = "exec"
)

func (k *Kubelet) DoProbe(pType ProbeType, containers map[string]containerd.Container, podConfig core.Pod) error {
	pStat := k.PodMap[podConfig.MetaData.Name]
	ready := true
	for id, container := range containers {
		//name := k.IDtoName[id]
		cStat := pStat.Containers[id]
		if pType == ExecProbe {
			// update pod status
			execOk := k.execProbe(container, podConfig, &cStat)
			pStat.Containers[id] = cStat
			if !execOk {
				pStat.Status = core.PhaseFailed
				k.PodMap[podConfig.MetaData.Name] = pStat
				// return an error
				return errors.New(constants.ErrorRestartPod)
			}
			if cStat.State.State != core.PhaseSucceed && !cStat.Ready {
				ready = false
			}
		}
		// currently do not support other kinds of probe
	}
	if ready == true {
		pStat.Status = core.PhaseRunning
		pStat.Conditions.Ready = true
	} else {
		pStat.Status = core.PhasePending
		pStat.Conditions.Ready = false
	}
	k.PodMap[podConfig.MetaData.Name] = pStat
	return nil
}

// execProbe 更新container status, 返回是否需要处理
func (k *Kubelet) execProbe(container containerd.Container, pConfig core.Pod, cStat *core.ContainerStatus) bool {
	retryTime := 5
	for i := 0; i < retryTime; i++ {
		status, err := utils.GetContainerStatus(container)
		if err != nil || status.Status == containerd.Unknown {
			continue
		}
		cStat.State.ExitCode = status.ExitStatus
		cStat.State.Finished = status.ExitTime
		if status.Status == containerd.Stopped {
			if status.ExitStatus == 0 {
				cStat.State.State = core.PhaseSucceed
				return true
			} else {
				cStat.State.State = core.PhaseFailed
				return pConfig.Spec.RestartPolicy == core.RestartNever
			}
		} else {
			if status.Status != containerd.Running {
				cStat.State.State = core.PhasePending
			} else {
				cStat.Ready = true
				cStat.State.State = core.PhaseRunning
			}
			return true
		}
	}
	cStat.State.State = core.PhaseUnknown
	return false
}
