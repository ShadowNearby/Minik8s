package kubelet

import (
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/go-connections/nat"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
)

func PodTest() {
	// create a pod config
	metadata := core.MetaData{
		Name:      "test",
		NameSpace: namespaces.Default,
		UUID:      "abcde",
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
	err := RunPod(&podConfig)
	if err != nil {
		logger.Errorf("run pod error: %s", err.Error())
	}
	StopPod(&podConfig)
}
