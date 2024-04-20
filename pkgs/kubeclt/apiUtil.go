package kubeclt

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	netrequest "minik8s/utils"
)

func GetApiKindFromYamlFile(fileContent []byte) (string, error) {
	var result map[string]interface{}

	err := yaml.Unmarshal(fileContent, &result)
	if err != nil {
		log.Debug("Kubectl", "GetApiKindFromYamlFile: Unmarshal object failed "+err.Error())
		return "", err
	}

	if result["kind"] == nil {
		log.Error("no kind found in file")
		return "", err
	}
	log.Infoln(result)
	return result["kind"].(string), nil
}
func ParseApiObjectFromYamlFile(fileContent []byte, obj interface{}) error {
	log.Debugln(fileContent)
	err := yaml.Unmarshal(fileContent, obj)

	if err != nil {
		log.Debug("Kubectl", "GetApiKindObjectFromYamlFile: Unmarshal object failed "+err.Error())
		return err
	}

	log.Debugln(obj)
	return err
}

func PostApiObjectToServer(URL string, obj interface{}) (int, error, string) {
	code, res, err := netrequest.PostRequestByTarget(URL, obj)
	if err != nil {
		log.Error("Kubectl", "PostApiObjectToServe: Post failed "+err.Error())
		return code, err, ""
	}
	bodyBytes, err := json.Marshal(res)
	if err != nil {
		return code, err, ""
	}
	return code, nil, string(bodyBytes)
}
func DeleteAPIObjectToServer(URL string) (int, error) {
	log.Debug("DeleteAPIObjectToServer", "URL: "+URL)
	code, err := netrequest.DelRequest(URL)
	if err != nil {
		log.Error("Kubectl", "DeleteAPIObjectToServer: Delete object failed "+err.Error())
		return code, err
	}
	log.Infoln("code: ", code)
	return code, nil
}
