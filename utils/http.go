package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var Resources = []string{"pod", "service", "endpoint", "replica", "job", "hpa", "dnsrecord"}
var Globals = []string{"function", "workflow", "node"}

//
//func ParseUrlOne(kind string, name string, ns string) string {
//	// operation: get. eg: "/api/v1/namespaces/{namespace}/pod/{pod_name}"
//	kind = strings.ToLower(kind)
//	name = strings.ToLower(name)
//	var namespace string
//	if ns == "nil" {
//		url := fmt.Sprintf("http://%s/api/v1/%s/%s", config.GetMasterIp(), kind, name)
//		return url
//	}
//	if ns == "" {
//		namespace = "default"
//	} else {
//		namespace = ns
//	}
//	url := fmt.Sprintf("http://%s/api/v1/namespaces/%s/%s/%s", config.GetMasterIp(), namespace, kind, name)
//	return url
//}
//func ParseUrlMany(kind string, ns string) string {
//	// operation: get. eg: GET "/api/v1/namespaces/{namespace}/pods"
//	// operation: create/apply. eg: POST "/api/v1/namespaces/{namespace}/pods"
//	var namespace string
//	if ns == "nil" {
//		url := fmt.Sprintf("http://%s/api/v1/%ss", config.GetMasterIp(), kind)
//		return url
//	}
//	if ns == "" {
//		namespace = "default"
//	} else {
//		namespace = ns
//	}
//	url := fmt.Sprintf("http://%s/api/v1/namespaces/%s/%ss", config.GetMasterIp(), namespace, kind)
//	return url
//}
//
////func ParseUrlFromJson(_json []byte) string {
////	// operation: create/apply. eg: POST "/api/v1/namespaces/{namespace}/pods"
////	kind := strings.ToLower(gjson.Get(string(_json), "kind").String())
////	namespace := gjson.Get(string(_json), "metadata.namespace")
////
////	url := fmt.Sprintf("http://%s/api/v1/namespaces/%s/%ss", config.ApiServerIp, namespace, kind)
////	return url
////}

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
