package replicaset

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

func TestReplicasetBasic(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logrus.SetReportCaller(true)

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "replicaset.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	rs := core.ReplicaSet{}
	err = json.Unmarshal(content, &rs)
	if err != nil {
		t.Errorf("parse repliaset error")
	}
	err = utils.CreateObject(core.ObjReplicaSet, "default", rs)
	if err != nil {
		t.Errorf("create rs failed: %s", err.Error())
	}
	time.Sleep(5 * time.Second)
	err = utils.DeleteObject(core.ObjReplicaSet, "default", rs.MetaData.Name)
	if err != nil {
		t.Errorf("delete rs failed: %s", err.Error())
	}
}

func TestReplicaUpdate(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logrus.SetReportCaller(true)

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "replicaset.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	rs := core.ReplicaSet{}
	err = json.Unmarshal(content, &rs)
	if err != nil {
		t.Errorf("parse repliaset error")
	}
	err = utils.CreateObject(core.ObjReplicaSet, "default", rs)
	if err != nil {
		t.Errorf("create rs failed: %s", err.Error())
	}
	time.Sleep(5 * time.Second)
	rs.Spec.Replicas += 1
	err = utils.SetObject(core.ObjReplicaSet, "default", rs.MetaData.Name, rs)
	if err != nil {
		t.Errorf("update rs failed: %s", err.Error())
	}

	time.Sleep(1 * time.Second)
	rs.Spec.Replicas -= 3
	err = utils.SetObject(core.ObjReplicaSet, "default", rs.MetaData.Name, rs)
	if err != nil {
		t.Errorf("update rs failed: %s", err.Error())
	}

	time.Sleep(5 * time.Second)
	err = utils.DeleteObject(core.ObjReplicaSet, "default", rs.MetaData.Name)
	if err != nil {
		t.Errorf("delete rs failed: %s", err.Error())
	}
}
