package api

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wxnacy/wgo/arrays"
	"gopkg.in/yaml.v3"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"strings"
)

func GetObjTypeFromYamlFile(fileContent []byte) (core.ObjType, error) {
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
	if idx := arrays.ContainsString(core.ObjTypeAll, strings.ToLower(result["kind"].(string)+"s")); idx == -1 {
		return "", errors.New("Error kind: " + result["kind"].(string))
	} else {
		if result["kind"].(string) == "ReplicaSet" {
			return core.ObjReplicaSet, err
		}

		return core.ObjType(strings.ToLower(result["kind"].(string)) + "s"), err
	}
}

func GetNameFromParamsFile(fileContent []byte) (string, error) {
	var result map[string]interface{}
	err := yaml.Unmarshal(fileContent, &result)
	if err != nil {
		log.Debug("Kubectl", "GetNameFromParamsFile: Unmarshal object failed "+err.Error())
		return "", err
	}
	return result["name"].(string), nil
}
func GetParamsFromParamsFile(fileContent []byte) (string, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal(fileContent, &data)
	if err != nil {
		return "", err
	}
	log.Infoln(data)
	params, ok := data["params"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("params not found in YAML")
	}
	// 将params转换为JSON
	jsonBytes := utils.JsonMarshal(params)
	return jsonBytes, nil
}

func ParseApiObjectFromYamlFile(fileContent []byte, obj interface{}) error {
	log.Debugln(fileContent)
	err := yaml.Unmarshal(fileContent, &obj)
	if err != nil {
		log.Debug("Kubectl", "GetApiKindObjectFromYamlFile: Unmarshal object failed "+err.Error())
		return err
	}
	log.Debugln(obj)
	return err
}
