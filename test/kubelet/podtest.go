package test

import (
	"bufio"
	"fmt"
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/go-connections/nat"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet"
	"minik8s/utils"
	"os"
)

func PodBasicTest() {
	// create a pod config
	metadata := core.MetaData{
		Name:      "test",
		NameSpace: namespaces.Default,
		UUID:      utils.GenerateUUID(),
	}
	portMap := nat.PortMap{}
	portMap[nat.Port(80)] = make([]nat.PortBinding, 0)
	portMap[nat.Port(80)] = append(portMap[nat.Port(80)], nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "9898",
	})
	exposedPorts := make([]string, 1)
	exposedPorts[0] = "80"
	containerConfig := core.Container{
		Name:            "nginx",
		Image:           "docker.io/library/nginx:latest",
		ImagePullPolicy: core.PullIfNeeds,
		Cmd:             nil,
		Args:            nil,
		WorkingDir:      "/home/nginx",
		VolumeMounts:    nil,
		PortBindings:    portMap,
		ExposedPorts:    exposedPorts,
		Env:             nil,
		Resources:       core.ResourcesConfig{},
	}
	specs := make([]core.Container, 1)
	specs[0] = containerConfig

	podConfig := core.Pod{
		ApiVersion: "v1",
		MetaData:   metadata,
		Spec: core.Spec{
			Containers:    specs,
			RestartPolicy: core.RestartOnFailure,
			NodeSelector:  nil,
		},
		Status: core.Status{},
	}
	err := kubelet.CreatePod(&podConfig)
	if err != nil {
		logger.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")

	res1, err := utils.NerdTest("ps")
	logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)

	kubelet.StopPod(&podConfig)
}

func PodLocalhostTest() {
	metadata := core.MetaData{
		Name:      "test",
		NameSpace: namespaces.Default,
		UUID:      utils.GenerateUUID(),
	}
	portMap := nat.PortMap{}
	portMap[nat.Port(80)] = make([]nat.PortBinding, 0)
	portMap[nat.Port(80)] = append(portMap[nat.Port(80)], nat.PortBinding{
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
		Image:           "docker.io/library/python:3.7-alpine",
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
		Image:           "docker.io/library/python:3.7-alpine",
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
		Spec: core.Spec{
			Containers:    specs,
			RestartPolicy: core.RestartOnFailure,
			NodeSelector:  nil,
		},
		Status: core.Status{},
	}
	err := kubelet.CreatePod(&podConfig)
	if err != nil {
		logger.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")
	res1, err := utils.NerdTest("ps")
	logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)
	fmt.Println("input c for terminate")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText := scanner.Text()
		if len(inputText) > 0 && inputText[0] == 'c' {
			break
		}
		if len(inputText) > 0 && inputText[0] == 'd' {
			res2, _ := utils.NerdTest("ps", "-a")
			res1, _ := utils.NerdTest("ps")
			logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)
		}
		fmt.Println("input c for terminate")
	}
	kubelet.StopPod(&podConfig)
}
