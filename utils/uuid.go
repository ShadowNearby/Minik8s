package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GenerateUUID generate a random uuid
func GenerateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		logrus.Errorf("generate uuid failed: %s", err.Error())
		return ""
	}
	return strings.Replace(id.String(), "-", "", -1)[:12]
}

// GenerateContainerIDByName receive container name + pod name as param
func GenerateContainerIDByName(containerName string, podUUID string) string {
	id := uuid.NewMD5(uuid.NameSpaceDNS, []byte(containerName+podUUID))
	return strings.Replace(id.String(), "-", "", -1)[:12]
}
