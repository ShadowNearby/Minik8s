package utils

import (
	"context"
	"errors"
	"fmt"
	core "minik8s/pkgs/apiobject"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	logger "github.com/sirupsen/logrus"
)

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

func GenerateLinuxNamespace(linuxNs map[string]string) []specs.LinuxNamespace {
	var namespaces []specs.LinuxNamespace
	for type_, item := range linuxNs {
		namespaces = append(namespaces, specs.LinuxNamespace{
			Type: specs.LinuxNamespaceType(type_),
			Path: item,
		})
	}
	return namespaces
}

func StopPodContainers(cStatus []core.ContainerStatus, pConfig core.Pod) error {
	var cs = make([]string, len(cStatus))
	for i, status := range cStatus {
		cs[i] = status.ID
	}
	_, err := NerdContainerOps(cs, pConfig.MetaData.Namespace, NerdStop)
	return err
}

func RmPodContainers(cStatus []core.ContainerStatus, pod core.Pod) error {
	var cs = make([]string, len(cStatus))
	for i, status := range cStatus {
		cs[i] = status.ID
	}
	_, err := NerdContainerOps(cs, pod.MetaData.Namespace, NerdRm)
	return err
}

func NerdContainerOps(containers []string, namespace string, ctlType string, args ...string) (string, error) {
	if ctlType == NerdInspect && len(containers) > 1 {
		return "", errors.New("cannot inspect no more than one container")
	}
	var retOutput string
	for _, c := range containers {
		output, err := NerdExec(NerdCtl{namespace: namespace, containerName: c, ctlType: ctlType}, args...)
		if err != nil {
			logger.Errorf("%s container %s failed: %s", ctlType, c, err.Error())
		}
		retOutput = output
	}
	return retOutput, nil
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
