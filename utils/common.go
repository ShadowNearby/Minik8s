package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"net"

	logger "github.com/sirupsen/logrus"
)

func JsonMarshal(item any) string {
	jsonText, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		logger.Errorf("marshal error: %s", err.Error())
	}
	return string(jsonText)
}

func JsonUnMarshal(text string, bind any) error {
	bytes := []byte(text)
	err := json.Unmarshal(bytes, bind)
	if err != nil {
		return err
	}
	return nil
}

func GetIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		if iface.Name != "eth0" && iface.Name != "ens3" {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
				continue
			}
			return ipNet.IP.String()
		}
	}
	return ""
}

func GenerateNewClusterIP() string {
	return fmt.Sprintf("%s.%d.%d", constants.IPPrefix, rand.Uint32()%256, rand.Uint32()%256)
}

func MatchLabel(l map[string]string, r map[string]string) bool {
	for k, v := range l {
		if val, ok := r[k]; ok != true || val != v {
			return false
		}
	}
	return true
}
func GetPodListFromRS(rs *core.ReplicaSet) []*core.Pod {
	var podList []*core.Pod
	info := []byte(GetObject(core.ObjPod, rs.MetaData.Namespace, ""))
	json.Unmarshal(info, &podList)
	return podList
}
