package runtime

import (
	"context"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/resources"
	"minik8s/utils"

	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
)

type ProbeType string

const (
	ExecProbe ProbeType = "exec"
)

func findContainerById(pod *core.Pod, ID string) *core.ContainerStatus {
	for _, status := range pod.Status.ContainersStatus {
		if len(ID) > 12 && ID[:12] == status.ID {
			return &status
		}
		if status.ID == ID {
			return &status
		}
	}
	return nil
}

func restartContainer(pod *core.Pod, container *containerd.Container, cStatus *core.ContainerStatus) (err error) {
	retryTimes := 5
	for i := 0; i < retryTimes; i++ {
		err = resources.ContainerManagerInstance.StartContainer(context.Background(), *container, pod)
		if err != nil {
			continue
		}
		cStatus.RestartCount += 1
		cStatus.State = core.ContainerRunning
		cStatus.TerminateState = ""
		break
	}
	return
}

func (k *Kubelet) DoProbe(pType ProbeType, containers []containerd.Container, pod *core.Pod) error {
	containerStatusList := make([]core.ContainerStatus, 0)
	for _, container := range containers {
		containerStatus := findContainerById(pod, container.ID())
		if containerStatus == nil {
			// remove the container
			logrus.Infof("killing container %s", container.ID())
			list := make([]core.ContainerStatus, 1)
			list[0] = core.ContainerStatus{ID: container.ID()}
			_ = utils.StopPodContainers(list, *pod)
			_ = utils.RmPodContainers(list, *pod)
			continue
		}
		if pType == ExecProbe {
			execOk := k.execProbe(container, containerStatus)
			if !execOk {
				// try to restart container
				if pod.Spec.RestartPolicy != core.RestartNever {
					err := restartContainer(pod, &container, containerStatus)
					if err != nil {
						// cannot restart
						logrus.Warnf("container %s restart failed", container.ID())
					} else {
						logrus.Infof("container %s restart success", container.ID())
					}
				}
			}
		}
		containerStatusList = append(containerStatusList, *containerStatus)
	}
	pod.Status.ContainersStatus = containerStatusList
	pod.Status.Phase = UpdatePodPhase(containerStatusList)
	return nil
}

func UpdatePodPhase(containerStatusList []core.ContainerStatus) core.PodPhase {
	allCount := len(containerStatusList)
	waitingCount := 0
	runningCount := 0
	terminatedCount := 0
	failCount := 0
	successCount := 0
	for _, status := range containerStatusList {
		switch status.State {
		case core.ContainerRunning:
			runningCount++
		case core.ContainerTerminated:
			if status.TerminateState == core.ContainerTerminateSuccess {
				successCount++
			}
			if status.TerminateState == core.ContainerTerminateFail {
				failCount++
			}
			terminatedCount++
		case core.ContainerWaiting:
			waitingCount++
		}
	}
	if waitingCount != 0 {
		return core.PodPhasePending
	}
	if runningCount == allCount {
		return core.PodPhaseRunning
	}
	if terminatedCount == (allCount-1) && successCount == (allCount-1) {
		return core.PodPhaseSucceeded
	}
	if terminatedCount == (allCount-1) && failCount > 0 {
		return core.PodPhaseFailed
	}
	return core.PodPhaseUnknown
}

// execProbe 更新container status, 返回是否ok
func (k *Kubelet) execProbe(container containerd.Container, cStatus *core.ContainerStatus) bool {
	retryTime := 5
	for i := 0; i < retryTime; i++ {
		status, err := utils.GetContainerStatus(container)
		if err != nil || status.Status == containerd.Unknown {
			continue
		}
		if status.Status == containerd.Stopped {
			// needs further check
			cStatus.State = core.ContainerTerminated

			if status.ExitStatus != 0 {
				cStatus.TerminateState = core.ContainerTerminateFail
				return true
			} else {
				cStatus.TerminateState = core.ContainerTerminateSuccess
				return false
			}
		} else {
			if status.Status == containerd.Running {
				cStatus.State = core.ContainerRunning
			} else {
				cStatus.State = core.ContainerWaiting
			}
			return true
		}
	}
	cStatus.State = core.ContainerWaiting
	return false
}
