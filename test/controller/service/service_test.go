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

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "pods.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	pods := []core.Pod{}
	err = json.Unmarshal(content, &pods)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	content, err = os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "services.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	services := []core.Service{}
	err = json.Unmarshal(content, &services)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}

	err = utils.CreateObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0])
	if err != nil {
		t.Errorf("error in create pod err: %s", err.Error())
	}
	err = utils.CreateObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1])
	if err != nil {
		t.Errorf("error in create pod err: %s", err.Error())
	}
	time.Sleep(10 * time.Second)

	err = utils.CreateObject(core.ObjService, services[0].MetaData.Namespace, services[0])
	if err != nil {
		t.Errorf("error in create service err: %s", err.Error())
	}
	time.Sleep(2 * time.Second)
	port := services[0].Spec.Ports[0].Port
	code, raw, err := utils.SendRequest("GET", fmt.Sprintf("http://%s:%d", services[0].Spec.ClusterIP, port), []byte{})
	if err != nil {
		t.Errorf("Error sending request: %s", err.Error())
	}

	if code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, code)
	}

	resp := string(raw)
	expect := "hello\n"
	if resp != expect {
		t.Errorf("Expected response %s, got %s", expect, resp)
	}

	err = utils.DeleteObject(core.ObjService, services[0].MetaData.Namespace, services[0].MetaData.Name)
	if err != nil {
		t.Errorf("error in delete service err: %s", err.Error())
	}

	err = utils.DeleteObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0].MetaData.Name)
	if err != nil {
		t.Errorf("error in delete pod err: %s", err.Error())
	}
	err = utils.DeleteObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1].MetaData.Name)
	if err != nil {
		t.Errorf("error in delete pod err: %s", err.Error())
	}

	time.Sleep(10 * time.Second)
}

func TestServiceWithDNS(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logrus.SetReportCaller(true)

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "pods.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	pods := []core.Pod{}
	err = json.Unmarshal(content, &pods)
	if err != nil {
		t.Errorf("Error unmarshalling pods: %s", err.Error())
	}
	content, err = os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "services.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	services := []core.Service{}
	err = json.Unmarshal(content, &services)
	if err != nil {
		t.Errorf("Error unmarshalling services: %s", err.Error())
	}

	err = utils.CreateObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0])
	if err != nil {
		t.Errorf("error in create pod err: %s", err.Error())
	}
	err = utils.CreateObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1])
	if err != nil {
		t.Errorf("error in create pod err: %s", err.Error())
	}
	time.Sleep(10 * time.Second)

	err = utils.CreateObject(core.ObjService, services[0].MetaData.Namespace, services[0])
	if err != nil {
		t.Errorf("error in create service err: %s", err.Error())
	}
	time.Sleep(2 * time.Second)

	content, err = os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "service_dns.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	dnsRecord := core.DNSRecord{}
	err = json.Unmarshal(content, &dnsRecord)
	if err != nil {
		t.Errorf("Error unmarshalling dns record: %s", err.Error())
	}
	err = utils.CreateObject(core.ObjDNS, dnsRecord.MetaData.Namespace, dnsRecord)
	if err != nil {
		t.Errorf("error in create dns err: %s", err.Error())
	}
	time.Sleep(2 * time.Second)

	code, raw, err := utils.SendRequest("GET", fmt.Sprintf("http://%s:%d/%s", dnsRecord.Host, 80, dnsRecord.Paths[0].Path), []byte{})
	if err != nil {
		t.Errorf("Error sending request: %s", err.Error())
	}

	if code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, code)
	}

	resp := string(raw)
	expect := "hello\n"
	if resp != expect {
		t.Errorf("Expected response %s, got %s", expect, resp)
	}

	err = utils.DeleteObject(core.ObjDNS, dnsRecord.MetaData.Namespace, dnsRecord.MetaData.Name)
	if err != nil {
		t.Errorf("error in del dns err: %s", err.Error())
	}

	err = utils.DeleteObject(core.ObjService, services[0].MetaData.Namespace, services[0].MetaData.Name)
	if err != nil {
		t.Errorf("error in del service err: %s", err.Error())
	}

	err = utils.DeleteObject(core.ObjPod, pods[0].MetaData.Namespace, pods[0].MetaData.Name)
	if err != nil {
		t.Errorf("error in del pod err: %s", err.Error())
	}
	err = utils.DeleteObject(core.ObjPod, pods[1].MetaData.Namespace, pods[1].MetaData.Name)
	if err != nil {
		t.Errorf("error in del pod err: %s", err.Error())
	}

	time.Sleep(10 * time.Second)
}
