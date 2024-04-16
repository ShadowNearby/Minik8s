package utils

import (
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
)

func GenerateUUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		logger.Errorf("generate uuid failed: %s", err.Error())
		return "", err
	}
	return id.String(), err
}

func GenerateContainerIDByName(str string) string {
	id := uuid.NewMD5(uuid.NameSpaceDNS, []byte(str))
	return id.String()[:12]
}
