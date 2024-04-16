package kubelet

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

func RunPod(podConfig *core.Pod) error {
	// first create a pause c_config
	var pauseConfig = core.Container{
		Name:            "pause_container",
		Image:           "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9",
		Cmd:             []string{"/pause"},
		ImagePullPolicy: core.PullIfNeeds,
	}
	var pauseSpec = utils.GenerateContainerSpec(*podConfig, pauseConfig)
	ctx := context.Background()
	var cMng ContainerManager
	var startedContainer = make([]core.Container, 0)
	_, err := cMng.CreateContainer(ctx, pauseSpec)
	if err != nil {
		logger.Errorf("Create Pause Container Failed: %s", err.Error())
		return err
	}
	startedContainer = append(startedContainer, pauseConfig)

	// get pause container namespace
	pausePid, err := cMng.GetContainerInfo(podConfig.MetaData.NameSpace, pauseSpec.Name, "State", "Pid")
	if err != nil {
		logger.Errorf("cannot get inspect: %s", err.Error())
	}
	linuxNamespace := fmt.Sprintf("/proc/%s/ns/", pausePid)

	// create pod containers
	for _, cConfig := range podConfig.Spec.Containers {
		// while create containers, add into pause container's namespace
		_, err := cMng.CreateContainer(ctx, utils.GenerateContainerSpec(*podConfig, cConfig, linuxNamespace))
		if err != nil {
			logger.Errorf("Create container %s Failed: %s", cConfig.Name, err.Error())
			_ = utils.StopPodContainers(startedContainer, podConfig.MetaData.NameSpace)
			_ = utils.RmPodContainers(startedContainer, podConfig.MetaData.NameSpace)
			return err
		}
		startedContainer = append(startedContainer, cConfig)
	}

	// add pause config to pod containers
	podConfig.Spec.Containers = append(podConfig.Spec.Containers, pauseConfig)
	return nil
}

func StopPod(podConfig *core.Pod) error {
	// stop every containerd in pod
	_ = utils.StopPodContainers(podConfig.Spec.Containers, podConfig.MetaData.NameSpace)
	_ = utils.RmPodContainers(podConfig.Spec.Containers, podConfig.MetaData.NameSpace)
	return nil
}
