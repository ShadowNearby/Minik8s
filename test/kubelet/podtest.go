package test

import (
	"bufio"
	"fmt"
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/go-connections/nat"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/resources"
	"minik8s/utils"
	"os"
	"testing"
	"time"
)

func PodBasicTest(t *testing.T) {
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
		Spec: core.Spec{
			Containers:    specs,
			RestartPolicy: core.RestartOnFailure,
			NodeSelector:  nil,
		},
		Status: core.Status{},
	}
	err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		logger.Errorf("run pod error: %s", err.Error())
		t.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")

	res1, err := utils.NerdTest("ps")
	logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)

	err = resources.StopPod(podConfig)
	if err != nil {
		t.Errorf("stop pod error: %s", err.Error())
	}
}

func PodLocalhostTest() {
	podConfig := GeneratePodConfigPy()
	err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		logger.Errorf("run pod error: %s", err.Error())
		//t.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")
	res1, err := utils.NerdTest("ps")
	logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)
	fmt.Println("input c for terminate")
	time.Sleep(5 * time.Second)
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
	err = resources.StopPod(podConfig)
	if err != nil {
		//t.Errorf("stop pod error: %s", err.Error())
	}
}
