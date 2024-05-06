package controller

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/config"
	"minik8s/utils"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func SetObject(objType core.ObjType, namespace string, name string, obj any) error {
	if namespace == "" {
		namespace = "default"
	}
	var url string
	objTxt := utils.JsonMarshal(obj)
	if name == "" {
		url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s", config.LocalServerIp, config.ApiServerPort, namespace, string(objType))

	} else {
		url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s/%s", config.LocalServerIp, config.ApiServerPort, namespace, string(objType), name)
	}
	if code, info, err := utils.SendRequest("PUT", url, []byte(objTxt)); err != nil || code != http.StatusOK {
		logger.Errorf("[set obj error]: %s", info)
		return err
	}
	return nil
}

func GetObject(objType core.ObjType, namespace string, name string) string {
	if namespace == "" {
		namespace = "default"
	}
	var url string
	if name == "" {
		url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s",
			config.LocalServerIp, config.ApiServerPort, namespace, string(objType))
	} else {
		url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s/%s",
			config.LocalServerIp, config.ApiServerPort, namespace, string(objType), name)
	}

	if code, info, err := utils.SendRequest("GET", url, make([]byte, 0)); err != nil || code != http.StatusOK {
		logger.Error("[get obj error]: ", info)
		return ""
	} else {
		return info
	}
}

func CreateObject(objType core.ObjType, namespace string, object any) error {
	if namespace == "" {
		namespace = "default"
	}
	var url string
	objectTxt := utils.JsonMarshal(object)
	url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s",
		config.LocalServerIp, config.ApiServerPort, namespace, objType)
	if code, info, err := utils.SendRequest("POST", url, []byte(objectTxt)); err != nil || code != http.StatusOK {
		logger.Errorf("[create obj error]: %s", info)
		return err
	} else {
		return nil
	}
}

func DeleteObject(objType core.ObjType, namespace string, name string) error {
	if namespace == "" {
		namespace = "default"
	}
	var url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s/%s",
		config.LocalServerIp, config.ApiServerPort, namespace, objType, name)
	if code, info, err := utils.SendRequest("DELETE", url, make([]byte, 0)); err != nil || code != http.StatusOK {
		logger.Errorf("[delete object error]: %s", info)
		return err
	} else {
		return nil
	}
}

func SplitChannelInfo(key string) (namespace, name string, err error) {
	parts := strings.Split(key, "/")
	switch len(parts) {
	case 1:
		// name only, no namespace
		return "", parts[0], nil
	case 2:
		// namespace and name
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected key format: %q", key)
}
