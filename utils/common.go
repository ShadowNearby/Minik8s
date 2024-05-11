package utils

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"net"
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
		if iface.Name != "eth0" {
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
