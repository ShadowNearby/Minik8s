package utils

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
)

func GenerateContainerName(pConfig core.Pod, cConfig core.Container) string {
	return fmt.Sprintf("%s_%s", cConfig.Name, pConfig.MetaData.NameSpace)
}

func GenerateContainerSpec(pConfig core.Pod, cConfig core.Container) core.ContainerdSpec {
	var cSpec = core.ContainerdSpec{
		Namespace:    pConfig.MetaData.NameSpace,
		Image:        cConfig.Image,
		Name:         GenerateContainerName(pConfig, cConfig),
		VolumeMounts: generateVolMountsMap(cConfig.VolumeMounts),
		Cmd:          cConfig.Cmd,
		Envs:         generateEnvList(cConfig.Env),
		Resource:     cConfig.Resources.Limit,
		Labels:       pConfig.MetaData.Labels,
		PullPolicy:   cConfig.ImagePullPolicy,
	}
	return cSpec
}

func CreateClientWithNamespace(namespace string) (*containerd.Client, error) {
	client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace(namespace))
	if err != nil {
		logger.Errorf("create client failed: %s", err.Error())
		return nil, err
	}
	return client, err
}

func GenerateUserOpts(user string) []oci.SpecOpts {
	var opts []oci.SpecOpts
	if user != "" {
		opts = append(opts, oci.WithUser(user), withResetAdditionalGIDs(), oci.WithAdditionalGIDs(user))
	}
	return opts
}

func withResetAdditionalGIDs() oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		s.Process.User.AdditionalGids = nil
		return nil
	}
}

func GenerateMounts(mountMap map[string]string) []specs.Mount {
	var res []specs.Mount
	for des, src := range mountMap {
		res = append(res, specs.Mount{
			Destination: des,
			Source:      src,
			Type:        "bind",
			Options:     []string{"bind"},
		})
	}
	return res
}

func StopStartedContainers(containers []string, namespace string) error {
	for _, c := range containers {
		output, err := NerdExec(NerdCtl{namespace: namespace, containerName: c, ctlType: NerdStop})
		if err != nil {
			logger.Error("stop container %s failed: %s", c, err.Error())
		}
		output, err = CtrExec(Ctr{
			ctrType:       CtrSnapshot,
			ctrOp:         CtrRm,
			containerName: c,
			namespace:     namespace,
		})
		logger.Infof("stop container info: %s", output)
	}
	return nil
}

func generateVolMountsMap(configs []core.VolumeMountConfig) map[string]string {
	var res = make(map[string]string)
	for _, config := range configs {
		res[config.ContainerPath] = config.HostPath
	}
	return res
}

func generateEnvList(envs []core.EnvConfig) []string {
	var res []string
	for _, env := range envs {
		res = append(res, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	return res
}
