package test

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	kubeletcontroller "minik8s/pkgs/kubelet/controller"
	"minik8s/utils"
	"net/http"
	"testing"
	"time"

	"github.com/containerd/containerd/namespaces"
	"github.com/docker/go-connections/nat"
	logger "github.com/sirupsen/logrus"
)

func TestPodBasicTest(t *testing.T) {
	// create a pod config
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
	err := kubeletcontroller.CreatePod(&podConfig)
	if err != nil {
		_ = kubeletcontroller.StopPod(podConfig)
		t.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")
	res1, err := utils.NerdTest("ps")
	logger.Infof("ps output:\n%s\nps -a output:\n%s\n", res1, res2)
	_ = kubeletcontroller.StopPod(podConfig)
}

func TestPodLocalhostTest(t *testing.T) {
	podConfig := generateConfig()
	err := kubeletcontroller.CreatePod(&podConfig)
	if err != nil {
		t.Errorf("run pod error: %s", err.Error())
		//t.Errorf("run pod error: %s", err.Error())
	}
	res2, err := utils.NerdTest("ps", "-a")
	res1, err := utils.NerdTest("ps")
	t.Logf("ps output:\n%s\nps -a output:\n%s\n", res1, res2)
	time.Sleep(2 * time.Second)
	_ = kubeletcontroller.StopPod(podConfig)
}

func TestPodUpdate(t *testing.T) {
	pod := generateConfig()
	code, info, err := utils.SendRequest("POST", fmt.Sprintf("http://127.0.0.1:8090/api/v1/namespaces/%s/pods/", pod.MetaData.Namespace), []byte(utils.JsonMarshal(pod)))
	if err != nil {
		t.Error(err.Error())
	}
	if code != http.StatusOK {
		t.Error("return bad: ", info)
	}
	logger.Info("create pod return 200")
	time.Sleep(2 * time.Second)
	pod.MetaData.Labels = map[string]string{"test": "haha"}
	code, info, err = utils.SendRequest("PUT", fmt.Sprintf("http://127.0.0.1:8090/api/v1/namespaces/%s/pods/%s", pod.MetaData.Namespace, pod.MetaData.Name), []byte(utils.JsonMarshal(pod)))
	if err != nil {
		t.Error(err.Error())
	}
	if code != http.StatusOK {
		t.Error("return bad: ", info)
	}
	logger.Info("create pod return 200")
	time.Sleep(2 * time.Second)
	code, info, err = utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:8090/api/v1/namespaces/%s/pods/%s", pod.MetaData.Namespace, pod.MetaData.Name), nil)
	if err != nil {
		t.Error(err.Error())
	}
	if code != http.StatusOK {
		t.Error("return bad: ", info)
	}
	logger.Info("create pod return 200")
	_ = kubeletcontroller.StopPod(pod)
}

func generateConfig() core.Pod {
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
