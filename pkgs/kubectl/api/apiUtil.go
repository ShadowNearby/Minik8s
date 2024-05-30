package api

import (
	core "minik8s/pkgs/apiobject"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	kind := result["kind"].(string)
	kind = strings.ToLower(kind)
	var objType core.ObjType
	for _, ty := range core.ObjTypeAll {
		if !strings.Contains(ty, kind) {
			continue
		}
		objType = core.ObjType(ty)
	}
	return objType, nil
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
