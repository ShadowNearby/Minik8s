package kubelet

import (
	"context"
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
	var startedContainer []string
	_, err := cMng.CreateContainer(ctx, pauseSpec)
	if err != nil {
		logger.Errorf("Create Pause Container Failed: %s", err.Error())
		return err
	}
	startedContainer = append(startedContainer, utils.GenerateContainerName(*podConfig, pauseConfig))

	for _, cConfig := range podConfig.Spec.Containers {
		_, err := cMng.CreateContainer(ctx, utils.GenerateContainerSpec(*podConfig, cConfig))
		if err != nil {
			logger.Errorf("Create c_config %s Failed: %s", cConfig.Name, err.Error())
			utils.StopStartedContainers(startedContainer, podConfig.MetaData.NameSpace)
			return err
		}
		startedContainer = append(startedContainer, utils.GenerateContainerName(*podConfig, cConfig))
	}

	// add pause config to pod containers
	podConfig.Spec.Containers = append(podConfig.Spec.Containers, pauseConfig)
	return nil
}
