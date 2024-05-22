package resources

import (
	"context"
	"errors"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/volume"
	"minik8s/utils"

	logger "github.com/sirupsen/logrus"
)

type NameIdPair struct {
	Name string
	ID   string
}

// CreatePod create containers and start them
func CreatePod(podConfig *core.Pod, pStatusChan chan<- core.PodStatus, cStatusChan chan<- core.ContainerStatus, done chan<- bool) error {
	// first create a pause c_config
	pStat := utils.InitPodStatus(podConfig)
	if pStatusChan != nil {
		pStatusChan <- pStat
	}
	utils.CheckPodMetaData(podConfig)
	var pauseConfig = core.Container{
		Name:            core.PauseContainerName,
		Image:           "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9",
		Cmd:             []string{"/pause"},
		ImagePullPolicy: core.PullIfNeeds,
	}
	var pauseSpec = utils.GenerateContainerSpec(*podConfig, pauseConfig)
	ctx := context.Background()
	output, err := utils.NerdRun([]string{"run", "-d", "--name", pauseSpec.Name, "--namespace", pauseSpec.Namespace, "--net", "flannel", "--label", fmt.Sprintf("%s=%s", "name", pauseSpec.Name), "--label", fmt.Sprintf("%s=%s", constants.MiniK8SPod, pauseSpec.PodName), pauseSpec.Image}...)
	var startedContainer = make([]core.Container, 0)
	if err != nil {
		logger.Errorf("Run Pause Container Failed: %s\n output: %s", err.Error(), output)
		return err
	}
	pauseSpec.ID = output[:12]

	// change pod ip
	pStat.PodIP, _ = ContainerManagerInstance.GetContainerInfo(podConfig.MetaData.Namespace, pauseSpec.ID, "NetworkSettings", "IPAddress")
	pStatusChan <- pStat

	if cStatusChan != nil {
		cStatusChan <- core.ContainerStatus{
			ID:           pauseSpec.ID,
			Name:         pauseSpec.Name,
			Image:        pauseSpec.Image,
			State:        core.ContainerRunning,
			RestartCount: 0,
			Environment:  nil,
		}
	}
	startedContainer = append(startedContainer, pauseConfig)
	logger.Infof("------CREATE PAUSE CONTAINER OVER--------")

	// get pause container information
	// ns
	pausePid, err := ContainerManagerInstance.GetContainerInfo(podConfig.MetaData.Namespace, pauseSpec.ID, "State", "Pid")
	if err != nil {
		logger.Errorf("cannot get namespace: %s", err.Error())
		//_ = utils.RmPodContainers(startedContainer, *podConfig)
		return err
	}
	linuxNamespace := fmt.Sprintf("/proc/%s/ns/", pausePid)
	logger.Infof("namespace: %s", linuxNamespace)

	err = utils.NerdCopy(fmt.Sprintf("%s:%s", pauseSpec.ID, config.ContainerResolvPath), config.TempResolvPath, podConfig.MetaData.Namespace)
	if err != nil {
		logger.Errorf("copy pause resolv.conf to local failed: %s", err.Error())
	}
	err = utils.AddCoreDns()
	if err != nil {
		logger.Errorf("edit resolv.conf failed: %s", err.Error())
	}
	err = utils.NerdCopy(fmt.Sprintf("%s:%s", pauseSpec.ID, config.ContainerHostsPath), config.TempHostsPath, podConfig.MetaData.Namespace)
	if err != nil {
		logger.Errorf("copy pause hosts to local failed: %s", err.Error())
	}

	if len(podConfig.Spec.Volumes) > 0 {
		err = volume.HandleDynamicVolumes(podConfig.Spec.Volumes)
		if err != nil {
			logger.Errorf("error in dynamic create volume err: %s", err.Error())
		}
	}

	// create pod containers
	for _, cConfig := range podConfig.Spec.Containers {
		// while create containers, add into pause container's namespace
		volume.HandleMount(cConfig.VolumeMounts)

		cSpec := utils.GenerateContainerSpec(*podConfig, cConfig, linuxNamespace)
		container, err := ContainerManagerInstance.CreateContainer(ctx, cSpec)
		if err != nil {
			logger.Errorf("Create container %s Failed: %s", cConfig.Name, err.Error())
			//_ = utils.StopPodContainers(startedContainer, *podConfig)
			//_ = utils.RmPodContainers(startedContainer, *podConfig)
			return err
		}

		err = utils.NerdCopy(config.TempResolvPath, fmt.Sprintf("%s:%s", cSpec.ID, config.ContainerResolvPath), podConfig.MetaData.Namespace)
		if err != nil {
			logger.Errorf("copy local resolv.conf to container %s failed: %s", err.Error(), cSpec.Name)
		}
		err = utils.NerdCopy(config.TempHostsPath, fmt.Sprintf("%s:%s", cSpec.ID, config.ContainerHostsPath), podConfig.MetaData.Namespace)
		if err != nil {
			logger.Errorf("copy local hosts to container %s failed: %s", err.Error(), cSpec.Name)
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
	err := utils.StopPodContainers(podConfig.Status.ContainersStatus, podConfig)
	if err != nil {
		return err
	}
	err = utils.RmPodContainers(podConfig.Status.ContainersStatus, podConfig)
	if err != nil {
		return err
	}
	return nil
}

func GetPodContainersMetrics(pod *core.Pod) ([]core.ContainerMetrics, error) {
	containers := ContainerManagerInstance.GetPodContainers(pod)
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
	}
	return res, nil
}
