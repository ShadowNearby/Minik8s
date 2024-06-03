package runtime

import (
	"context"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/resources"
	"minik8s/utils"

	"github.com/containerd/containerd"
)

type ProbeType string

const (
	ExecProbe ProbeType = "exec"
)

func findContainerById(pod *core.Pod, ID string) *core.ContainerStatus {
	for _, status := range pod.Status.ContainersStatus {
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
			list := make([]core.ContainerStatus, 1)
			list[0] = core.ContainerStatus{ID: container.ID()}
			_ = utils.StopPodContainers(list, *pod)
			_ = utils.RmPodContainers(list, *pod)
			continue
		}
		if pType == ExecProbe {
			execOk := k.execProbe(container, pod, containerStatus)
			if !execOk {
				pod.Status.Condition = core.CondPending
				// try to restart container
				err := restartContainer(pod, &container, containerStatus)
				if err != nil {
					// cannot restart
					pod.Status.Condition = core.CondFailed
					return err
				}
			}
		}
		containerStatusList = append(containerStatusList, *containerStatus)
	}
	pod.Status.ContainersStatus = containerStatusList
	return nil
}

// execProbe 更新container status, 返回是否ok
func (k *Kubelet) execProbe(container containerd.Container, pod *core.Pod, cStatus *core.ContainerStatus) bool {
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
				return pod.Spec.RestartPolicy == core.RestartNever
			} else {
				return true
			}
		} else {
			if status.Status != containerd.Running {
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
