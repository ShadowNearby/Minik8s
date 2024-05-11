package utils

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"net/http"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func SetObject(objType core.ObjType, namespace string, name string, obj any, sendMsg ...bool) error {
	if namespace == "" {
		namespace = "default"
	}
	var url string
	objTxt := JsonMarshal(obj)
	if name == "" {
		url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s", config.LocalServerIp, config.ApiServerPort, namespace, string(objType))

	} else {
		url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s/%s", config.LocalServerIp, config.ApiServerPort, namespace, string(objType), name)
	}
	if len(sendMsg) == 0 || sendMsg[0] == true {
		if code, info, err := SendRequest("PUT", url, []byte(objTxt)); err != nil || code != http.StatusOK {
			logger.Errorf("[set obj error]: %s", info)
			return err
		}
	}
	return nil
}
func SetObjectWONamespace(objType core.ObjType, name string, obj any) error {
	var url string
	objTxt := JsonMarshal(obj)
	if name == "" {
		url = fmt.Sprintf("http://%s:%s/api/v1/%s", config.LocalServerIp, config.ApiServerPort, string(objType))

	} else {
		url = fmt.Sprintf("http://%s:%s/api/v1/%s/%s", config.LocalServerIp, config.ApiServerPort, string(objType), name)
	}
	if code, info, err := SendRequest("PUT", url, []byte(objTxt)); err != nil || code != http.StatusOK {
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

	var retInfo core.InfoType
	if code, info, err := SendRequest("GET", url, make([]byte, 0)); err != nil || code != http.StatusOK {
		_ = JsonUnMarshal(info, &retInfo)
		logger.Error("[get obj error]: ", retInfo.Error)
		return ""
	} else {
		_ = JsonUnMarshal(info, &retInfo)
		return retInfo.Data
	}
}
func GetObjectWONamespace(objType core.ObjType, name string) string {
	var url string
	if name == "" {
		url = fmt.Sprintf("http://%s:%s/api/v1/%s",
			config.LocalServerIp, config.ApiServerPort, string(objType))
	} else {
		url = fmt.Sprintf("http://%s:%s/api/v1/%s/%s",
			config.LocalServerIp, config.ApiServerPort, string(objType), name)
	}

	var retInfo core.InfoType
	if code, info, err := SendRequest("GET", url, make([]byte, 0)); err != nil || code != http.StatusOK {
		_ = JsonUnMarshal(info, &retInfo)
		logger.Error("[get obj error]: ", retInfo.Error)
		return ""
	} else {
		_ = JsonUnMarshal(info, &retInfo)
		return retInfo.Data
	}
}

func CreateObject(objType core.ObjType, namespace string, object any) error {
	if namespace == "" {
		namespace = "default"
	}
	var url string
	objectTxt := JsonMarshal(object)
	logger.Debugln(objectTxt)
	url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s",
		config.LocalServerIp, config.ApiServerPort, namespace, objType)
	if code, info, err := SendRequest("POST", url, []byte(objectTxt)); err != nil || code != http.StatusOK {
		logger.Errorf("[create obj error]: %s", info)
		return err
	} else {
		return nil
	}
}

func CreateObjectWONamespace(objType core.ObjType, object any) error {
	var url string
	objectTxt := JsonMarshal(object)
	logger.Debugln(objectTxt)
	url = fmt.Sprintf("http://%s:%s/api/v1/%s",
		config.LocalServerIp, config.ApiServerPort, objType)
	if code, info, err := SendRequest("POST", url, []byte(objectTxt)); err != nil || code != http.StatusOK {
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
	var url string
	url = fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/%s/%s",
		config.LocalServerIp, config.ApiServerPort, namespace, objType, name)
	if code, info, err := SendRequest("DELETE", url, make([]byte, 0)); err != nil || code != http.StatusOK {
		logger.Errorf("[delete object error]: %s", info)
		return err
	} else {
		return nil
	}
}

func DeleteObjectWONamespace(objType core.ObjType, name string) error {
	var url string
	url = fmt.Sprintf("http://%s:%s/api/v1/%s/%s",
		config.LocalServerIp, config.ApiServerPort, objType, name)
	if code, info, err := SendRequest("DELETE", url, make([]byte, 0)); err != nil || code != http.StatusOK {
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
