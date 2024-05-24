package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GenerateUUID generate a random uuid
func GenerateUUID(length ...int) string {
	id, err := uuid.NewUUID()
	if err != nil {
		logrus.Errorf("generate uuid failed: %s", err.Error())
		return ""
	}
	if len(length) == 0 {
		return strings.Replace(id.String(), "-", "", -1)[:12]
	} else {
		return strings.Replace(id.String(), "-", "", -1)[:length[0]]
	}
}

// GenerateContainerIDByName receive container name + pod name as param
func GenerateContainerIDByName(containerName string, podUUID string) string {
	id := uuid.NewMD5(uuid.NameSpaceDNS, []byte(containerName+podUUID))
	return strings.Replace(id.String(), "-", "", -1)[:12]
}
