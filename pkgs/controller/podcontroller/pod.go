package podcontroller

import (
	"encoding/json"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"

	log "github.com/sirupsen/logrus"
)

type PodController struct{}

func (pc *PodController) GetChannel() string {
	return constants.ChannelPod
}

func (pc *PodController) HandleCreate(message string) error {
	return nil
}

func (pc *PodController) HandleUpdate(message string) error {
	return nil
}

func (pc *PodController) HandleDelete(message string) error {
	pod := &core.Pod{}
	err := json.Unmarshal([]byte(message), pod)
	if err != nil {
		log.Errorf("unmarshal pod error: %s", err.Error())
		return err
	}
	log.Infof("delete pod: %s:%s", pod.MetaData.Namespace, pod.MetaData.Name)
	_, _, err = utils.SendRequest("DELETE", fmt.Sprintf("http://%s:%s/pod/stop/%s/%s", pod.Status.HostIP, config.NodePort, pod.GetNamespace(), pod.MetaData.Name), nil)
	if err != nil {
		log.Errorf("delete pod error: %s", err.Error())
		return err
	}
	return nil
}
