package test

import (
	"encoding/json"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"net/http"
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
	time.Sleep(10 * time.Second)

	utils.CreateObject(core.ObjService, services[0].MetaData.Namespace, services[0])
	time.Sleep(2 * time.Second)
	port := services[0].Spec.Ports[0].Port
	code, raw, err := utils.SendRequest("GET", fmt.Sprintf("http://%s:%d", services[0].Spec.ClusterIP, port), []byte{})
	if err != nil {
		logrus.Errorf("Error sending request: %s", err.Error())
	}

	if code != http.StatusOK {
		logrus.Errorf("Expected status code %d, got %d", http.StatusOK, code)
	}

	logrus.Infof("Response: %s", string(raw))

	utils.DeleteObject(core.ObjService, services[0].MetaData.Namespace, services[0].MetaData.Name)

	utils.DeleteObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0].MetaData.Name)
	utils.DeleteObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1].MetaData.Name)

	time.Sleep(10 * time.Second)
}
