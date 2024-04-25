package utils

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
)

func CreateJson(item any) string {
	jsonText, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		logger.Errorf("marshal error: %s", err.Error())
	}
	return string(jsonText)
}
