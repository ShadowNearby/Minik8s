package runtime

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"os"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

type Kubelet struct {
	MasterIP     string
	MasterPort   string
	Labels       map[string]string
	PodStatMap   map[string]core.PodStatus
	PodConfigMap map[string]core.Pod
	Server       *gin.Engine
	IDtoName     map[string]string
	NumCores     int
	MemCapacity  uint64
}

var KubeletInstance Kubelet

func (k *Kubelet) InitKubelet(config core.KubeletConfig) {
	k.MasterIP = config.MasterIP
	k.MasterPort = config.MasterPort
	k.Labels = config.Labels
	name, _ := os.Hostname()
	pods := []core.Pod{}
	_, resp, err := utils.SendRequest("GET", fmt.Sprintf("http://%s:%s/api/v1/nodes/%s/pods", k.MasterIP, k.MasterPort, name), []byte{})
	if err != nil {
		logger.Errorf("error in init pods in node")
	}
	respType := core.InfoType{}
	err = utils.JsonUnMarshal(resp, &respType)
	if err != nil {
		logger.Errorf("error in unmarshal resp in node")
	}
	err = utils.JsonUnMarshal(respType.Data, &pods)
	if err != nil {
		logger.Errorf("error in unmarshal pods in node")
	}
	k.PodStatMap = make(map[string]core.PodStatus)
	k.PodConfigMap = make(map[string]core.Pod)
	for _, pod := range pods {
		k.PodConfigMap[fmt.Sprintf("%s-%s", pod.MetaData.Name, pod.MetaData.Namespace)] = pod
		k.PodStatMap[fmt.Sprintf("%s-%s", pod.MetaData.Name, pod.MetaData.Namespace)] = pod.Status
	}
	k.Server = gin.Default()
}

func (k *Kubelet) RegisterNode() {
	name, _ := os.Hostname()
	nodeInfo := core.Node{
		ApiVersion: "v1",
		Kind:       "Node",
		NodeMetaData: core.MetaData{
			Name:   name,
			Labels: k.Labels,
		},
		Spec: core.NodeSpec{
			NodeIP:  utils.GetIP(),
			PodCIDR: config.PodCIDR,
		},
	}
	code, data, err := utils.SendRequest("POST",
		fmt.Sprintf("http://%s:%s/api/v1/nodes", k.MasterIP, k.MasterPort),
		[]byte(utils.JsonMarshal(nodeInfo)))
	if err != nil {
		logger.Errorf("send request error: %s", err.Error())
		return
	}
	if code != 200 {
		logger.Errorf("server error: %d, info: %s", code, data)
	}
}

func (k *Kubelet) GetPodConfig(podName string, podNamespace string) core.Pod {
	return k.PodConfigMap[fmt.Sprintf("%s-%s", podName, podNamespace)]
}

func (k *Kubelet) WritePodConfig(podName string, podNamespace string, podConfig *core.Pod) {
	if k.PodConfigMap == nil {
		k.PodConfigMap = make(map[string]core.Pod)
	}
	k.PodConfigMap[fmt.Sprintf("%s-%s", podName, podNamespace)] = *podConfig
}

func (k *Kubelet) DelPodConfig(podName string, podNamespace string) {
	if k.PodConfigMap == nil {
		return
	}
	delete(k.PodConfigMap, fmt.Sprintf("%s-%s", podName, podNamespace))
}
func (k *Kubelet) DelPodStat(podName string, podNamespace string) {
	if k.PodStatMap == nil {
		return
	}
	delete(k.PodStatMap, fmt.Sprintf("%s-%s", podName, podNamespace))
}
