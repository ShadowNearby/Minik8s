package cmd

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubectl/api"
	"minik8s/utils"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply can create apiObject in a declarative way",
	Long:  "Kubectl apply can create apiObject in a declarative way, usage kubectl apply -f [file]",
	Run:   applyHandler,
}

func applyHandler(cmd *cobra.Command, args []string) {
	// 检查参数是否是文件 读取文件
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	if fileInfo.IsDir() {
		log.Errorf("%s is not a file", filePath)
		return
	}
	// 读取文件的内容
	fileContent, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 解析API对象的种类
	objType, err := api.GetObjTypeFromYamlFile(fileContent)
	if err != nil {
		log.Fatal(err)
		return
	}
	structType, res := core.ObjTypeToCoreObjMap[objType]
	if !res {
		log.Error("Unsupported struct", objType)
		return
	}
	object := reflect.New(structType).Interface().(core.ApiObjectKind)
	err = yaml.Unmarshal(fileContent, object)
	if err != nil {
		log.Fatal(err)
	}
	nameSpace := object.GetNamespace()
	err = utils.CreateObject(objType, nameSpace, object)
	if err != nil {
		log.Fatal(err)
	}
}
