package scheduler

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"net/http"
)

type Scheduler struct {
	Policy     string `json:"policy"`
	podChannel <-chan *redis.Message
}

func (sched *Scheduler) Run(policy string) {
	sched.Policy = policy
	sched.podChannel = storage.RedisInstance.SubscribeChannel(constants.ChannelPodSchedule)
	go func() {
		for message := range sched.podChannel {
			msg := message.Payload
			pods := make([]core.Pod, 2)
			utils.JsonUnMarshal(msg, &pods)
			logger.Infof("old:%s\nnew:%s", utils.JsonMarshal(pods[0].Status), utils.JsonMarshal(pods[1].Status))
			pods[1].Status.HostIP = pods[0].Status.HostIP
			_, err := sched.Schedule(pods[1])
			if err != nil {
				logger.Errorf("schedule fail: %s", err.Error())
			}
		}
	}()
	select {}
	logger.Error("should not be here")
}

func (sched *Scheduler) Schedule(pod core.Pod) (string, error) {
	if pod.Status.HostIP != "" {
		err := sendStopPod(pod.Status.HostIP, pod)
		if err != nil {
			logger.Errorf("cannot stop pod on %s", pod.Status.HostIP)
		}
	}
	podSelector := pod.Spec.Selector.MatchLabels
	nodesTxt := utils.GetObjectWONamespace(core.ObjNode, "")
	var nodes []core.Node
	utils.JsonUnMarshal(nodesTxt, &nodes)
	var nodeCandidate = make(map[string]core.NodeMetrics)
	for _, node := range nodes {
		nodeLabels, metrics, err := requestNodeInfos(node)
		if err != nil {
			logger.Errorf("get node %s 's info failed", node.Spec.NodeIP)
			continue
		}
		flag := true
		for key, val := range podSelector {
			if v, ok := nodeLabels[key]; ok != true || v != val {
				flag = false
				break
			}
		}
		if flag {
			nodeCandidate[node.Spec.NodeIP] = metrics
		}
	}
	if len(nodeCandidate) == 0 {
		return "", errors.New("cannot schedule the pod")
	}
	selectedIP := sched.dispatch(nodeCandidate)
	pod.Status.HostIP = selectedIP
	// select node over, send message to node
	err := sendCreatePod(selectedIP, pod)
	if err != nil {
		return "", err
	}
	// node register pod over, write back to storage
	err = utils.SetObject(core.ObjPod, pod.MetaData.Namespace, pod.MetaData.Name, pod)
	if err != nil {
		return "", err
	}
	return selectedIP, nil
}

func (sched *Scheduler) dispatch(candidates map[string]core.NodeMetrics) string {
	switch sched.Policy {
	case config.PolicyCPU:
		{
			return cpuPolicy(candidates)
		}
	case config.PolicyMemory:
		{
			return memPolicy(candidates)
		}
	case config.PolicyDisk:
		{
			return diskPolicy(candidates)
		}
	default:
		{
			logger.Errorf("unsuppported policy")
			return ""
		}

	}
}

func requestNodeInfos(node core.Node) (map[string]string, core.NodeMetrics, error) {
	//url := fmt.Sprintf("http://%s:%s/metrics", node.Spec.NodeIP, config.NodePort)
	// TODO: using ip
	url := fmt.Sprintf("http://%s:%s/metrics", "127.0.0.1", config.NodePort)
	code, data, err := utils.SendRequest("GET", url, []byte(""))
	if err != nil || code != http.StatusOK {
		logger.Error("get metrics error")
		return nil, core.NodeMetrics{}, err
	}
	var info core.InfoType
	var metrics core.NodeMetrics
	utils.JsonUnMarshal(data, &info)
	utils.JsonUnMarshal(info.Data, &metrics)
	return node.NodeMetaData.Labels, metrics, nil
}

func sendCreatePod(nodeIp string, pod core.Pod) error {
	//url := fmt.Sprintf("http://%s:%s/pod/create", nodeIp, config.NodePort)
	// TODO: using ip
	url := fmt.Sprintf("http://%s:%s/pod/create", "127.0.0.1", config.NodePort)
	code, info, err := utils.SendRequest("POST", url, []byte(utils.JsonMarshal(pod)))
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		var data core.InfoType
		utils.JsonUnMarshal(info, &data)
		return errors.New(data.Error)
	}
	return nil
}

func sendStopPod(nodeIP string, pod core.Pod) error {
	//url := fmt.Sprintf("http://%s:%s/pod/stop",nodeIP, config.NodePort)
	// TODO: using ip
	url := fmt.Sprintf("http://%s:%s/pod/stop", "127.0.0.1", config.NodePort)
	code, info, err := utils.SendRequest("POST", url, []byte(utils.JsonMarshal(pod)))
	if err != nil {
		return nil
	}
	if code != http.StatusOK {
		var data core.InfoType
		utils.JsonUnMarshal(info, &data)
		return errors.New(data.Error)
	}
	return nil
}