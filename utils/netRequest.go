package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func DelRequest(uri string) (int, error) {
	status_code, resp_buffer, err := SendRequest(http.MethodDelete, uri, nil)
	log.Debug(resp_buffer)
	return status_code, err
}

func PostRequestByTarget(uri string, target interface{}) (int, interface{}, error) {
	jsonData, err := json.Marshal(target)
	if err != nil {
		log.Error("postRequest", "PostRequestByTarget: Marshal object failed "+err.Error())
		return 0, nil, err
	}

	status_code, resp, err := SendRequestWithJson(http.MethodPost, "application/json", jsonData)
	var bodyJson interface{}
	if err := json.Unmarshal([]byte(resp), &bodyJson); err != nil {
		log.Error("postRequest", "PostRequestByTarget: Decode response failed "+err.Error())
		return 0, nil, err
	}
	return status_code, bodyJson, nil

}
