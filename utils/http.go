package utils

import (
	"bytes"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func SendRequest(method string, url string, body []byte) (string, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
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
		log.Infoln(response.Status)
		log.Debugln("[Http Request] to %s method:%s status:%s receive:%d bytes", url, method, response.Status, length)
	}

	return buffer.String(), err
}

func SendRequestWithJson(method string, url string, json []byte) (string, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		return "", err
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
		log.Debugln("[Http Request] to %s method:%s status:%s receive:%d bytes", url, method, response.Status, length)
	}
	return buffer.String(), err
}

func SendRequestWithHost(method string, url string, body []byte) (string, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
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
	return buffer.String(), err
}
