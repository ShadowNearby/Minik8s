package test

import (
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/go-connections/nat"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

func GeneratePodConfigPy() core.Pod {
	metadata := core.MetaData{
		Name:      "test",
		Namespace: namespaces.Default,
		UUID:      utils.GenerateUUID(),
	}
	portMap := nat.PortMap{}
	portMap[nat.Port(rune(80))] = make([]nat.PortBinding, 0)
	portMap[nat.Port(rune(80))] = append(portMap[nat.Port(rune(80))], nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "9898",
	})
	exposedPorts := make([]string, 1)
	exposedPorts[0] = "80"
	volumeMount := core.VolumeMountConfig{
		ContainerPath: "/home/python",
		HostPath:      "../../../../../home/k8s/ly/minik8s/test/kubelet",
		ReadOnly:      false,
	}
	env1 := core.EnvConfig{
		Name:  "PORT_SERVER",
		Value: "8080",
	}
	env4 := core.EnvConfig{
		Name:  "PORT_CLIENT",
		Value: "8080",
	}
	containerConfig := core.Container{
		Name:            "py1",
		Image:           "docker.io/library/python:3.9-alpine",
		ImagePullPolicy: core.PullIfNeeds,
		Cmd:             []string{"python3", "/home/python/server.py"},
		Args:            nil,
		WorkingDir:      "/home/nginx",
		VolumeMounts:    []core.VolumeMountConfig{volumeMount},
		PortBindings:    portMap,
		ExposedPorts:    exposedPorts,
		Env:             []core.EnvConfig{env1},
		Resources:       core.ResourcesConfig{},
	}
	specs := make([]core.Container, 2)
	specs[0] = containerConfig
	containerConfig2 := core.Container{
		Name:            "py2",
		Image:           "docker.io/library/python:3.9-alpine",
		ImagePullPolicy: core.PullIfNeeds,
		Cmd:             []string{"python3", "/home/python/client.py"},
		Args:            nil,
		WorkingDir:      "/home/nginx",
		VolumeMounts:    []core.VolumeMountConfig{volumeMount},
		PortBindings:    portMap,
		ExposedPorts:    exposedPorts,
		Env:             []core.EnvConfig{env4},
		Resources:       core.ResourcesConfig{},
	}
	specs[1] = containerConfig2
	podConfig := core.Pod{
		ApiVersion: "v1",
		MetaData:   metadata,
		Spec: core.PodSpec{
			Containers:    specs,
			RestartPolicy: core.RestartOnFailure,
			Selector:      core.Selector{},
		},
		Status: core.PodStatus{},
	}
	return podConfig
}
