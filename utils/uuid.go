package utils

import (
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
)

// GenerateUUID generate a random uuid
func GenerateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		logger.Errorf("generate uuid failed: %s", err.Error())
		return ""
	}
	return id.String()
}

// GenerateContainerIDByName receive container name + pod name as param
func GenerateContainerIDByName(containerName string, podName string) string {
	id := uuid.NewMD5(uuid.NameSpaceDNS, []byte(containerName+podName))
	return id.String()[:12]
}
