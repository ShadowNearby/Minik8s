package test

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestPrometheusPod(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableQuote: true, ForceColors: true})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)

	pod := &core.Pod{}
	path := fmt.Sprintf("%s/%s", utils.ExamplePath, "prometheus_pod.json")
	content, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("error in read file %s err: %s", path, err.Error())
		return
	}
	err = utils.JsonUnMarshal(string(content), pod)
	if err != nil {
		logrus.Errorf("error in unmarshal err: %s", err.Error())
		return
	}
	err = utils.CreateObject(core.ObjPod, pod.MetaData.Namespace, pod)
	if err != nil {
		t.Errorf("create pod error: %s", err.Error())
	}

	time.Sleep(30 * time.Second)
	err = utils.DeleteObject(core.ObjPod, pod.MetaData.Namespace, pod.MetaData.Name)
	if err != nil {
		t.Errorf("del pod error: %s", err.Error())
	}
}