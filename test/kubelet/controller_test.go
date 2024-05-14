package test

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"testing"
	"time"
)

func TestInspectPod(t *testing.T) {
	podConfig := utils.GeneratePodConfigPy()
	err := kubeletcontroller.CreatePod(&podConfig)
	if err != nil {
		_ = kubeletcontroller.StopPod(podConfig)
		t.Errorf("cannot create pod")
		return
	}
	kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)

	time.Sleep(6 * time.Second)

	kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)

	_ = kubeletcontroller.StopPod(podConfig)
}

func TestMetricsTest(t *testing.T) {
	podConfig := utils.GeneratePodConfigPy()
	err := kubeletcontroller.CreatePod(&podConfig)
	//err := resources.CreatePod(&podConfig, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot create pod")
		kubeletcontroller.StopPod(podConfig)
		return
	}
	_ = kubeletcontroller.NodeMetrics()
	//text := utils.JsonMarshal(stat)
	kubeletcontroller.StopPod(podConfig)
}

func TestHttpInspectPod(t *testing.T) {
	pod := utils.GeneratePodConfigPy()
	code, _, err := utils.SendRequest("POST", "http://127.0.0.1:10250/pod/create", []byte(utils.JsonMarshal(pod)))
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:10250/pod/stop/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
		t.Error("create pod error")
		return
	}
	code, data, err := utils.SendRequest("GET", fmt.Sprintf("http://127.0.0.1:10250/pod/status/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:10250/pod/stop/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
		t.Error("inspect pod error")
		return
	}
	var info core.InfoType
	utils.JsonUnMarshal(data, &info)
	var inspect core.PodStatus
	err = utils.JsonUnMarshal(info.Data, &inspect)
	if err != nil {
		t.Error("unmarshal inspection error")
		return
	}
	//logger.Info(utils.JsonMarshal(inspect))
	utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:10250/pod/stop/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
}

func TestHttpPodMetrics(t *testing.T) {
	pod := utils.GeneratePodConfigPy()
	logger.Info(len(pod.Spec.Containers))
	code, _, err := utils.SendRequest("POST", "http://127.0.0.1:10250/pod/create", []byte(utils.JsonMarshal(pod)))
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:10250/pod/stop/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
		t.Error("create pod error")
		return
	}
	code, data, err := utils.SendRequest("GET", fmt.Sprintf("http://127.0.0.1:10250/metrics/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:10250/pod/stop/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
		t.Error("find metrics error")
		return
	}
	var info core.InfoType
	utils.JsonUnMarshal(data, &info)
	var metrics []core.ContainerMetrics
	utils.JsonUnMarshal(info.Data, &metrics)
	logger.Info(utils.JsonMarshal(metrics))
	utils.SendRequest("DELETE", fmt.Sprintf("http://127.0.0.1:10250/pod/stop/%s/%s", pod.GetNameSpace(), pod.MetaData.Name), nil)
}
