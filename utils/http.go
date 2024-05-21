package utils

import (
	"bytes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Resources = []string{"pod", "service", "endpoint", "replica", "job", "hpa", "dnsrecord"}
var Globals = []string{"function", "workflow", "node"}

func ParseJson(c *gin.Context) map[string]any {
	json := make(map[string]any)
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return make(map[string]any)
	}
	return json
}

func SendRequest(method string, url string, body []byte) (int, string, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, "", err
	}
	request.Header.Set("content-type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	buffer := &bytes.Buffer{}
	if err != nil {
		log.Error(err)
		if response == nil {
			return http.StatusInternalServerError, "", err
		}
	} else {
		length, err := buffer.ReadFrom(response.Body)
		if err != nil {
			log.Error(err)
		}
		err = response.Body.Close()
		if err != nil {
			log.Error(err)
		}
		// log.Infoln(response.Status)
		log.Infof("[Send Request] to %s method:%s status:%s receive:%d bytes", url, method, response.Status, length)
	}

	return response.StatusCode, buffer.String(), err
}

func SendRequestWithJson(method string, url string, json []byte) (int, string, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		return 0, "", err
	}
	request.Header.Set("content-type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	buffer := &bytes.Buffer{}
	if err != nil {
		log.Error(err)
	} else {
		length, err := buffer.ReadFrom(response.Body)
		if err != nil {
			log.Error(err)
		}
		err = response.Body.Close()
		if err != nil {
			log.Error(err)
		}
		log.Debugf("[Http Request] to %s method:%s status:%s receive:%d bytes", url, method, response.Status, length)
	}
	return response.StatusCode, buffer.String(), err
}

func SendRequestWithHost(method string, url string, body []byte) (int, string, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, "", err
	}
	request.Header.Set("content-type", "application/json")
	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err)
	}
	request.Header.Set("source", hostname)
	client := &http.Client{}
	response, err := client.Do(request)
	buffer := &bytes.Buffer{}
	if err != nil {
		log.Error(err)
	} else {
		length, err := buffer.ReadFrom(response.Body)
		if err != nil {
			log.Error(err)
		}
		err = response.Body.Close()
		if err != nil {
			log.Error(err)
		}
		log.Debugf("[Http Request] to %s method:%s status:%s receive:%d bytes", url, method, response.Status, length)
	}
	return response.StatusCode, buffer.String(), err
}
