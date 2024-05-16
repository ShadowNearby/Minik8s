package test

import (
	"encoding/json"
	core "minik8s/pkgs/apiobject"
	kubeletcontroller "minik8s/pkgs/kubelet/controller"
	"minik8s/utils"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestServiceController(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logrus.SetReportCaller(true)

	content, err := os.ReadFile("pods.json")
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	pods := []core.Pod{}
	json.Unmarshal(content, &pods)
	content, err = os.ReadFile("services.json")
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	services := []core.Service{}
	json.Unmarshal(content, &services)
	logrus.Infof("%d\n", len(services))

	utils.CreateObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0])
	utils.CreateObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1])
	time.Sleep(5 * time.Second)
	utils.DeleteObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0].MetaData.Name)
	utils.DeleteObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1].MetaData.Name)
	_ = kubeletcontroller.StopPod(pods[0])
	_ = kubeletcontroller.StopPod(pods[1])
}
