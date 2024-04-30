package service

import (
	"fmt"
	core "minik8s/pkgs/apiobject"

	log "github.com/sirupsen/logrus"
)

func FindUnusedIP() string {
	for i, used := range UsedIP {
		if i == 0 || used {
			continue
		}
		return fmt.Sprintf("%s%d", IPPrefix, i)
	}
	log.Errorf("No IP available")
	return ""
}

func MatchLabel(l map[string]string, r map[string]string) bool {
	for k, v := range l {
		if r[k] != v {
			return false
		}
	}
	return true
}

func FindDestPort(targetPort string, containers []core.Container) uint32 {
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.Name == targetPort {
				return p.ContainerPort
			}
		}
	}
	return 0
}
