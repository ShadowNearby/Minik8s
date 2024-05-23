package hpa

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"os"
	"testing"
	"time"
)

func TestHpaBasic(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logrus.SetReportCaller(true)

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "replicaset.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	rs := core.ReplicaSet{}
	err = json.Unmarshal(content, &rs)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	content, err = os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "hpa.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	hpa := core.HorizontalPodAutoscaler{}
	err = json.Unmarshal(content, &hpa)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}

	// create replicaset first
	err = utils.CreateObject(core.ObjReplicaSet, "default", rs)
	if err != nil {
		t.Errorf("error in create rs: %s", err.Error())
	}
	logrus.Infof("create success")
	// then create hpa
	err = utils.CreateObject(core.ObjHpa, "default", hpa)
	if err != nil {
		t.Errorf("error in create hpa: %s", err.Error())
	}
	logrus.Infof("create success")
}
func TestHpaUpdate(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logrus.SetReportCaller(true)

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "replicaset.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	rs := core.ReplicaSet{}
	err = json.Unmarshal(content, &rs)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	content, err = os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "hpa.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	hpa := core.HorizontalPodAutoscaler{}
	err = json.Unmarshal(content, &hpa)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}

	// create replicaset first
	err = utils.CreateObject(core.ObjReplicaSet, "default", rs)
	if err != nil {
		t.Errorf("error in create rs: %s", err.Error())
	}
	// then create hpa
	err = utils.CreateObject(core.ObjHpa, "default", hpa)
	if err != nil {
		t.Errorf("error in create hpa: %s", err.Error())
	}
	time.Sleep(5 * time.Second)
	hpa.Spec.MaxReplicas = 5
	err = utils.SetObject(core.ObjHpa, "default", hpa.MetaData.Name, hpa)
	if err != nil {
		t.Errorf("error in update hpa: %s", err.Error())
	}

}
