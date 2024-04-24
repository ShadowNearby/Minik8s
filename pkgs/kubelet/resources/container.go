package resources

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"strings"
)

type ContainerManager struct {
	Client *containerd.Client
}

var ContainerManagerInstance ContainerManager

func (cmg *ContainerManager) createClient(namespace string) {
	if cmg.Client == nil {
		client, err := utils.CreateClientWithNamespace(namespace)
		if err != nil {
			logger.Errorf("Create Client Error: %s", err.Error())
			panic(err)
		}
		logger.Info("Create Containerd Client Success!")
		cmg.Client = client
	}
}

func (cmg *ContainerManager) CreateContainer(ctx context.Context, config core.ContainerdSpec) (containerd.Container, error) {
	cmg.createClient(config.Namespace)
	var imgCtl ImageController
	img, err := imgCtl.CreateImage(cmg.Client, config.Image, config.PullPolicy) // TODO: safer
	if err != nil {
		logger.Errorf("Create Image %s Failed: %s", config.Image, err.Error())
		return nil, err
	}
	logger.Infof("Create/Find Image %s Success!", config.Image)

	// define specs and options
	specs := []oci.SpecOpts{oci.WithImageConfig(img)}
	// add user and group information
	//specs = append(specs, utils.GenerateUserOpts("k8s")...)
	if len(config.VolumeMounts) > 0 {
		specs = append(specs, oci.WithMounts(utils.GenerateMounts(config.VolumeMounts)))
	}
	if len(config.Cmd) > 0 {
		specs = append(specs, oci.WithProcessArgs(config.Cmd...))
	}
	if config.Resource.Cpu != core.EmptyCpu {
		specs = append(specs, oci.WithCPUs(config.Resource.Cpu))
	}
	if config.Resource.Memory != core.EmptyMemory {
		specs = append(specs, oci.WithMemoryLimit(config.Resource.Memory))
	}
	if len(config.Envs) > 0 {
		specs = append(specs, oci.WithEnv(config.Envs))
	}
	// add network support
	if len(config.LinuxNamespace) > 0 {
		linuxNamespaces := utils.GenerateLinuxNamespace(config.LinuxNamespace)
		for _, namespace := range linuxNamespaces {
			specs = append(specs, oci.WithLinuxNamespace(namespace))
		}
	}
	copts := []containerd.NewContainerOpts{containerd.WithImageName(config.Name), containerd.WithNewSnapshot(config.Name, img), containerd.WithNewSpec(specs...)}
	if len(config.Labels) > 0 {
		copts = append(copts, containerd.WithContainerLabels(config.Labels))
	}
	// add filter labels
	copts = append(copts, containerd.WithAdditionalContainerLabels(utils.GenerateContainerLabel(config.PodName)))
	// create container
	container, err := cmg.Client.NewContainer(ctx, config.ID, copts...)
	if err != nil {
		logger.Errorf("Create Container Failed: %s", err.Error())
		return nil, err
	}
	logger.Infof("Create Container %s Success", config.Name)
	return container, nil
}

//func (cmg ContainerManager) StartContainer(ctx context.Context, container containerd.Container) error {
//
//}

func (cmg *ContainerManager) StartContainer(ctx context.Context, container containerd.Container, pConfig *core.Pod) error {
	cmg.createClient(pConfig.MetaData.NameSpace)
	task, err := container.NewTask(ctx, cio.NewCreator())
	if err != nil {
		logger.Errorf("create task error: %s", err.Error())
		return err
	}
	err = task.Start(ctx)
	if err != nil {
		logger.Errorf("start tast error: %s", err.Error())
		return err
	}
	return nil
}

func (cmg *ContainerManager) GetContainerInfo(namespace string, containerID string, fields ...string) (string, error) {
	var str = ""
	for _, field := range fields {
		str += "." + field
	}
	str = fmt.Sprintf("{{%s}}", str)
	res, err := utils.NerdContainerOps([]string{containerID}, namespace, utils.NerdInspect, "-f", str)
	if err != nil {
		logger.Errorf("inspect error: %s", err.Error())
	}
	return strings.Trim(res, "\n "), err
}

func (cmg *ContainerManager) GetPodContainers(pConfig *core.Pod) []containerd.Container {
	cmg.createClient(pConfig.MetaData.NameSpace)
	cs, err := cmg.Client.Containers(context.Background(), fmt.Sprintf("labels.%q==%s", constants.MiniK8SPod, pConfig.MetaData.Name))
	if err != nil {
		logger.Errorf("filter containers failed: %s", err.Error())
		return nil
	}
	return cs
}
