package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
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

func StopPodContainers(containers []core.Container, namespace string) error {
	var cs = make([]string, len(containers))
	for i, container := range containers {
		cs[i] = fmt.Sprintf("%s_%s", container.Name, namespace)
	}
	_, err := NerdContainerOps(cs, namespace, NerdStop)
	return err
}

func RmPodContainers(containers []core.Container, namespace string) error {
	var cs = make([]string, len(containers))
	for i, container := range containers {
		//cs[i] = fmt.Sprintf("%s_%s", container.Name, namespace)
		cs[i] = container.Name
	}
	// rm container
	_, _ = NerdContainerOps(cs, namespace, NerdRm)
	//_ = ctlContainerOps(cs, namespace, CtrSnapshot, CtrRm)
	return nil
}

func ctlContainerOps(containers []string, namespace string, ctrObject string, ctrType string) error {
	for _, c := range containers {
		output, err := CtrExec(Ctr{
			ctrType:       ctrObject,
			ctrOp:         ctrType,
			containerName: c,
			namespace:     namespace,
		})
		if err != nil {
			logger.Errorf("rm snapshot error: %s", err.Error())
		}
		logger.Infof("rm snapshot info: %s", string(output))
	}
	return nil
}

func NerdContainerOps(containers []string, namespace string, ctlType string, args ...string) (string, error) {
	if ctlType == NerdInspect && len(containers) > 1 {
		return "", errors.New("cannot inspect no more than one container")
	}
	var retOutput string
	for _, c := range containers {
		output, err := NerdExec(NerdCtl{namespace: namespace, containerName: c, ctlType: ctlType}, args...)
		if err != nil {
			logger.Error("stop container %s failed: %s", c, err.Error())
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
