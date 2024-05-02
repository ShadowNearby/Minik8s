package api

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wxnacy/wgo/arrays"
	"gopkg.in/yaml.v3"
	core "minik8s/pkgs/apiobject"
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
	if idx := arrays.ContainsString(core.ObjTypeAll, strings.ToLower(result["kind"].(string))); idx == -1 {
		return "", errors.New("Error kind: " + result["kind"].(string))
	} else {
		if result["kind"].(string) == "ReplicaSet" {
			return core.ObjReplicaSet, err
		}

		return core.ObjType(strings.ToLower(result["kind"].(string)) + "s"), err
	}
}

func GetCoreObjFromObjType(objType core.ObjType) (interface{}, bool) {
	obj, exists := core.ObjTypeToCoreObjMap[objType]
	return obj, exists
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
