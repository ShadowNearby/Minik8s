package test

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"os"
	"testing"
)

func TestServiceController(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

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

	utils.DeleteObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0].MetaData.Name)
	utils.DeleteObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1].MetaData.Name)
}
