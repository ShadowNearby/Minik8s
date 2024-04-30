package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/kubectl"
	"minik8s/utils"
	"os"
	"reflect"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply can create apiObject in a declarative way",
	Long:  "Kubectl apply can create apiObject in a declarative way, usage kubectl apply -f [file]",
	Run:   applyHandler,
}

func applyHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("kubectl apply f [file] is required")
		return
	}
	// 检查参数是否是文件 读取文件
	fileInfo, err := os.Stat(args[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	if fileInfo.IsDir() {
		log.Errorf("%s is not a file", args[0])
		return
	}
	// 读取文件的内容
	fileContent, err := utils.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	// 解析API对象的种类
	objType, err := kubectl.GetObjTypeFromYamlFile(fileContent)
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
	nameSpace := object.GetNameSpace()
	log.Debugln(object)
	err = controller.CreateObject(objType, nameSpace, object)
	if err != nil {
		log.Fatal(err)
	}
}
