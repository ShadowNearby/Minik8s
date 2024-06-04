package runtime

import (
	"context"
	"fmt"
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
		cStatus.Status = containerd.Running
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
			execOk := k.execProbe(container, pod, containerStatus)
			if !execOk {
				// try to restart container
				if pod.Spec.RestartPolicy == core.RestartNever {
					logrus.Warnf("pod %s:%s failed", pod.MetaData.Namespace, pod.MetaData.Name)
					pod.Status.Condition = core.CondFailed
					return fmt.Errorf("pod failed")
				}
				err := restartContainer(pod, &container, containerStatus)
				if err != nil {
					// cannot restart
					logrus.Warnf("pod %s:%s failed", pod.MetaData.Namespace, pod.MetaData.Name)
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
		cStatus.Status = status.Status
		if status.Status == containerd.Stopped {
			// needs further check
			if status.ExitStatus != 0 {
				return false
			}
		}
		// running pausing paused created stopped(with exit 0)
		return true
	}
	cStatus.Status = containerd.Unknown
	return false
}
