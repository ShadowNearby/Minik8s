package utils

import (
	"encoding/json"
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
