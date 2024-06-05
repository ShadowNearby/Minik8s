package test

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"testing"

	logger "github.com/sirupsen/logrus"
)

func TestHttpInspectPod(t *testing.T) {
	ip := utils.GetIP()
	pod := utils.GeneratePodConfigPy()
	code, _, err := utils.SendRequest("POST", fmt.Sprintf("http://%s:10250/pod/create", ip), []byte(utils.JsonMarshal(pod)))
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://%s:10250/pod/stop/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
		t.Error("create pod error")
		return
	}
	code, data, err := utils.SendRequest("GET", fmt.Sprintf("http://%s:10250/pod/status/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://%s:10250/pod/stop/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
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
	utils.SendRequest("DELETE", fmt.Sprintf("http://%s:10250/pod/stop/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
}

func TestHttpPodMetrics(t *testing.T) {
	ip := utils.GetIP()
	pod := utils.GeneratePodConfigPy()
	logger.Info(len(pod.Spec.Containers))
	code, _, err := utils.SendRequest("POST", fmt.Sprintf("http://%s:10250/pod/create", ip), []byte(utils.JsonMarshal(pod)))
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://%s:10250/pod/stop/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
		t.Error("create pod error")
		return
	}
	code, data, err := utils.SendRequest("GET", fmt.Sprintf("http://%s:10250/metrics/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
	if code != 200 || err != nil {
		utils.SendRequest("DELETE", fmt.Sprintf("http://%s:10250/pod/stop/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
		t.Error("find metrics error")
		return
	}
	var info core.InfoType
	utils.JsonUnMarshal(data, &info)
	var metrics []core.ContainerMetrics
	utils.JsonUnMarshal(info.Data, &metrics)
	logger.Info(utils.JsonMarshal(metrics))
	utils.SendRequest("DELETE", fmt.Sprintf("http://%s:10250/pod/stop/%s/%s", ip, pod.GetNamespace(), pod.MetaData.Name), nil)
}
