package resources

import (
	"context"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

type NameIdPair struct {
	Name string
	ID   string
}

// CreatePod create containers and start them
func CreatePod(podConfig *core.Pod, pStatusChan chan<- core.PodStatus, cStatusChan chan<- core.ContainerStatus, done chan<- bool) error {
	// first create a pause c_config
	if pStatusChan != nil {
		pStatusChan <- utils.InitPodStatus(podConfig)
	}
	var pauseConfig = core.Container{
		Name:            core.PauseContainerName,
		Image:           "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9",
		Cmd:             []string{"/pause"},
		ImagePullPolicy: core.PullIfNeeds,
	}
	var pauseSpec = utils.GenerateContainerSpec(*podConfig, pauseConfig)
	ctx := context.Background()
	var startedContainer = make([]core.Container, 0)
	pause, err := ContainerManagerInstance.CreateContainer(ctx, pauseSpec)
	if err != nil {
		logger.Errorf("Create Pause Container Failed: %s", err.Error())
		return err
	}
	startedContainer = append(startedContainer, pauseConfig)
	err = ContainerManagerInstance.StartContainer(ctx, pause, podConfig)
	if err != nil {
		logger.Errorf("Start Pause Container Failed: %s", err.Error())
		//_ = utils.RmPodContainers(startedContainer, *podConfig)
		return err
	}
	logger.Infof("------CREATE PAUSE CONTAINER OVER--------")

	// get pause container information
	// ns
	pausePid, err := ContainerManagerInstance.GetContainerInfo(podConfig.MetaData.NameSpace, pauseSpec.ID, "State", "Pid")
	if err != nil {
		logger.Errorf("cannot get namespace: %s", err.Error())
		//_ = utils.RmPodContainers(startedContainer, *podConfig)
		return err
	}
	linuxNamespace := fmt.Sprintf("/proc/%s/ns/", pausePid)
	logger.Infof("namespace: %s", linuxNamespace)

	// create pod containers
	for _, cConfig := range podConfig.Spec.Containers {
		// while create containers, add into pause container's namespace
		container, err := ContainerManagerInstance.CreateContainer(ctx, utils.GenerateContainerSpec(*podConfig, cConfig, linuxNamespace))
		if err != nil {
			logger.Errorf("Create container %s Failed: %s", cConfig.Name, err.Error())
			//_ = utils.StopPodContainers(startedContainer, *podConfig)
			//_ = utils.RmPodContainers(startedContainer, *podConfig)
			return err
		}
		startedContainer = append(startedContainer, cConfig)
		err = ContainerManagerInstance.StartContainer(ctx, container, podConfig)
		if err != nil {
			logger.Errorf("Create container %s Failed: %s", cConfig.Name, err.Error())
			//_ = utils.StopPodContainers(startedContainer, *podConfig)
			//_ = utils.RmPodContainers(startedContainer, *podConfig)
			return err
		}
		if cStatusChan != nil {
			cStatusChan <- core.ContainerStatus{
				ID:           utils.GenerateContainerIDByName(cConfig.Name, podConfig.MetaData.UUID),
				Name:         cConfig.Name,
				Image:        cConfig.Image,
				State:        core.ContainerRunning,
				RestartCount: 0,
				Environment:  cConfig.Env,
			}
		}
		//runtime.KubeletInstance.ContainerStart(podConfig, cConfig.Name)
		logger.Infof("Start container %s Success", cConfig.Name)
	}
	// add pause config to pod containers
	podConfig.Spec.Containers = append(podConfig.Spec.Containers, pauseConfig)
	if done != nil {
		done <- true
	}
	return nil
}

// StopPod stop and remove the containers in pod
func StopPod(podConfig core.Pod) error {
	// stop every containerd in pod
	_ = utils.StopPodContainers(podConfig.Status.ContainersStatus, podConfig)
	_ = utils.RmPodContainers(podConfig.Status.ContainersStatus, podConfig)
	return nil
}

func GetPodMetrics(podConfig *core.Pod) ([]core.ContainerMetrics, error) {
	containers := ContainerManagerInstance.GetPodContainers(podConfig)
	if len(containers) == 0 {
		logger.Errorf("cannot filter pod's containers")
		return nil, errors.New("cannot filter container")
	}
	logger.Infof("----------------BEGIN TO INSPECT------------------")
	var res = make([]core.ContainerMetrics, 0)
	for _, container := range containers {
		metric, err := utils.GetContainerMetrics(container)
		if err != nil {
			logger.Errorf("inspect container %s failed: %s", container.ID(), err.Error())
			continue
		}
		res = append(res, metric)
		//logger.Infof("metric: %v", metric)
		//status, err := utils.GetContainerStatus(container)
		//if err != nil {
		//	logger.Errorf("inspect container %s failed: %s", container.ID(), err.Error())
		//	return err
		//}
		////status = status // change type to containerd.PodStatus
		//logger.Infof("status: %v", status)
	}
	return res, nil
}
