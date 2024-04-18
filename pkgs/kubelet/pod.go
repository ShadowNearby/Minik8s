package kubelet

import (
	"context"
	"errors"
	"fmt"
	"github.com/containerd/containerd"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"time"
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

	// get pause container information
	// ns
	pausePid, err := cMng.GetContainerInfo(podConfig.MetaData.NameSpace, pauseSpec.ID, "State", "Pid")
	if err != nil {
		logger.Errorf("cannot get inspect: %s", err.Error())
	}
	linuxNamespace := fmt.Sprintf("/proc/%s/ns/", pausePid)
	logger.Infof("namespace: %s", linuxNamespace)

	// stop pause container
	//_ = utils.StopPodContainers(startedContainer, podConfig)
	//startedContainer = make([]core.Container, 0)

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

func InspectPod(podConfig *core.Pod) error {
	var cmg ContainerManager
	containers := cmg.GetPodContainers(podConfig)
	if len(containers) == 0 {
		logger.Errorf("cannot filter pod's containers")
		return errors.New("cannot filter container")
	}
	logger.Infof("----------------BEGIN TO INSPECT------------------")
	var ranger []int
	ranger = append(ranger, 1, 2, 3)
	for _ = range ranger { // TODO: remove test
		time.Sleep(3 * time.Second)
		logger.Infof("wake up")
		for _, container := range containers {
			metric, err := runtime.GetContainerMetrics(container)
			if err != nil {
				logger.Errorf("inspect container %s failed: %s", container.ID(), err.Error())
				return err
			}
			logger.Infof("metric: %v", metric)
			status, err := runtime.GetContainerStatus(container)
			if err != nil {
				logger.Errorf("inspect container %s failed: %s", container.ID(), err.Error())
				return err
			}
			status = status.(containerd.Status) // change type to containerd.Status
			logger.Infof("status: %v", status)
		}
	}
	return nil
}
