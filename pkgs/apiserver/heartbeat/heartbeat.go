package heartbeat

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"time"

	"github.com/sirupsen/logrus"
)

func CheckNodeConnection() error {
	key := "/nodes/object"
	var nodes []core.Node
	err := storage.RangeGet(key, &nodes)
	if err != nil {
		logrus.Errorf("error in get node data")
		return err
	}
	for _, node := range nodes {
		lastHeartbeat := node.Status.LastHeartbeat
		if time.Since(lastHeartbeat) > config.HeartbeatInterval {
			logrus.Infof("node %s NetworkUnavailable since %s", node.NodeMetaData.Name, node.Status.LastHeartbeat.String())
			node.Status.Phase = core.NodeNetworkUnavailable
			nodeKey := fmt.Sprintf("%s/%s", key, node.NodeMetaData.Name)
			storage.Put(nodeKey, node)
		}
	}
	return nil
}

func Run() {
	for {
		time.Sleep(3 * config.HeartbeatInterval)
		err := CheckNodeConnection()
		if err != nil {
			logrus.Errorf("error in CheckNodeConnection %s", err.Error())
			return
		}
	}
}
