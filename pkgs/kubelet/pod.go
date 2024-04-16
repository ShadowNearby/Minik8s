package kubelet

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

// CreatePod create containers and start them
func CreatePod(podConfig *core.Pod) error {
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
	pause, err := cMng.CreateContainer(ctx, pauseSpec)
	if err != nil {
		logger.Errorf("Create Pause Container Failed: %s", err.Error())
		return err
	}
	startedContainer = append(startedContainer, pauseConfig)
	err = cMng.StartContainer(ctx, pause, podConfig)
	if err != nil {
		logger.Errorf("Start Pause Container Failed: %s", err.Error())
		_ = utils.RmPodContainers(startedContainer, podConfig)
		return err
	}
	logger.Infof("------CREATE PAUSE CONTAINER OVER--------")

	// get pause container namespace
	pausePid, err := cMng.GetContainerInfo(podConfig.MetaData.NameSpace, pauseSpec.ID, "State", "Pid")
	if err != nil {
		logger.Errorf("cannot get inspect: %s", err.Error())
	}
	linuxNamespace := fmt.Sprintf("/proc/%s/ns/", pausePid)
	logger.Infof("namespace: %s", linuxNamespace)

	// create pod containers
	for _, cConfig := range podConfig.Spec.Containers {
		// while create containers, add into pause container's namespace
		container, err := cMng.CreateContainer(ctx, utils.GenerateContainerSpec(*podConfig, cConfig, linuxNamespace))
		if err != nil {
			logger.Errorf("Create container %s Failed: %s", cConfig.Name, err.Error())
			_ = utils.StopPodContainers(startedContainer, podConfig)
			_ = utils.RmPodContainers(startedContainer, podConfig)
			return err
		}
		startedContainer = append(startedContainer, cConfig)
		err = cMng.StartContainer(ctx, container, podConfig)
		if err != nil {
			logger.Errorf("Create container %s Failed: %s", cConfig.Name, err.Error())
			_ = utils.StopPodContainers(startedContainer, podConfig)
			_ = utils.RmPodContainers(startedContainer, podConfig)
			return err
		}
		logger.Infof("Start container %s Success", cConfig.Name)
	}
	// add pause config to pod containers
	podConfig.Spec.Containers = append(podConfig.Spec.Containers, pauseConfig)
	return nil
}

// StopPod stop and remove the containers in pod
func StopPod(podConfig *core.Pod) error {
	// stop every containerd in pod
	_ = utils.StopPodContainers(podConfig.Spec.Containers, podConfig)
	_ = utils.RmPodContainers(podConfig.Spec.Containers, podConfig)
	return nil
}
