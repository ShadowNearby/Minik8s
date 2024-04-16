package kubelet

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

type ContainerManager struct {
	Client *containerd.Client
}

func (cmg *ContainerManager) CreateContainer(ctx context.Context, config core.ContainerdSpec) (containerd.Container, error) {
	if cmg.Client == nil {
		client, err := utils.CreateClientWithNamespace(config.Namespace)
		if err != nil {
			logger.Errorf("Create Client Error: %s", err.Error())
			return nil, err
		}
		logger.Info("Create Containerd Client Success!")
		cmg.Client = client
	}
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
	copts := []containerd.NewContainerOpts{containerd.WithNewSnapshot(config.Name, img), containerd.WithNewSpec(specs...)}
	if len(config.Labels) > 0 {
		copts = append(copts, containerd.WithContainerLabels(config.Labels))
	}

	// create container
	container, err := cmg.Client.NewContainer(ctx, config.Name, copts...)
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

func (cmg *ContainerManager) GetContainerInfo(namespace string, containerName string, fields ...string) (string, error) {
	var str = ""
	for _, field := range fields {
		str += "." + field
	}
	str = fmt.Sprintf("{{%s}}", str)
	res, err := utils.NerdContainerOps([]string{containerName}, namespace, utils.NerdInspect, "-f", str)
	if err != nil {
		logger.Errorf("inspect error: %s", err.Error())
	}
	return res, err
}
