package test

import (
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/go-connections/nat"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/controller"
	"minik8s/utils"
	"testing"
	"time"
)

func TestPodBasicTest(t *testing.T) {
	// create a pod config
	metadata := core.MetaData{
		Name:      "test",
		NameSpace: namespaces.Default,
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
		Spec: core.PodSpec{
			Containers:    specs,
			RestartPolicy: core.RestartOnFailure,
			Selector:      core.Selector{},
		},
		Status: core.PodStatus{},
	}
	err := controller.CreatePod(&podConfig)
	if err != nil {
		_ = controller.StopPod(podConfig)
		t.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")
	res1, err := utils.NerdTest("ps")
	logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)

	_ = controller.StopPod(podConfig)
}

func TestPodLocalhostTest(t *testing.T) {
	podConfig := GeneratePodConfigPy()
	err := controller.CreatePod(&podConfig)
	if err != nil {
		t.Errorf("run pod error: %s", err.Error())
		//t.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")
	res1, err := utils.NerdTest("ps")
	t.Logf("ps output:\n%s\nps -a output:\n%s\n", res1, res2)
	time.Sleep(2 * time.Second)
	_ = controller.StopPod(podConfig)
}
