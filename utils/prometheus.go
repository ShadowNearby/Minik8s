package utils

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GeneratePrometheusNodeFile() error {
	nodes := []core.Node{}
	content := GetObjectWONamespace(core.ObjNode, "")
	err := JsonUnMarshal(content, &nodes)
	if err != nil {
		logrus.Errorf("error in unmarshal node %s", err.Error())
		return err
	}

	configs := []core.PrometheusSdConfig{}
	for _, node := range nodes {
		if node.Status.Phase == core.NodeNetworkUnavailable {
			continue
		}
		configs = append(configs, core.PrometheusSdConfig{
			Targets: []string{fmt.Sprintf("%s:%d", node.Spec.NodeIP, 8080), fmt.Sprintf("%s:%d", node.Spec.NodeIP, 9100)},
			Labels: map[string]string{
				"job":  "node",
				"host": fmt.Sprintf("%s:%d", node.Spec.NodeIP, 9100),
			},
		})
	}
	filePath := fmt.Sprintf("%s/%s", ConfigPath, config.PrometheusNodeFilePath)
	data := JsonMarshal(configs)
	err = os.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		logrus.Errorf("error in write file %s", err.Error())
		return err
	}
	return nil
}

func GeneratePrometheusPodFile() error {
	pods := []core.Pod{}
	content := GetObjectWONamespace(core.ObjPod, "")
	err := JsonUnMarshal(content, &pods)
	if err != nil {
		logrus.Errorf("error in unmarshal node %s", err.Error())
		return err
	}
	configs := []core.PrometheusSdConfig{}

	for _, pod := range pods {
		if pod.Status.Phase != core.PodPhaseRunning {
			continue
		}
		if pod.Status.PodIP == "" {
			continue
		}
		ip := pod.Status.PodIP
		podLabels := pod.MetaData.Labels
		port := 0
		for key, val := range podLabels {
			if key == constants.MiniK8SPrometheusPort {
				port, err = strconv.Atoi(val)
				if err != nil {
					logrus.Errorf("error in convert port %s", err.Error())
					continue
				}
			}
		}
		if port == 0 {
			continue
		}
		configs = append(configs, core.PrometheusSdConfig{
			Targets: []string{fmt.Sprintf("%s:%d", ip, port)},
			Labels: map[string]string{
				"job": "pod",
			},
		})
	}

	filePath := fmt.Sprintf("%s/%s", ConfigPath, config.PrometheusPodFilePath)
	err = os.WriteFile(filePath, []byte(JsonMarshal(configs)), 0644)
	if err != nil {
		logrus.Errorf("error in write file %s", err.Error())
		return err
	}
	return nil
}
