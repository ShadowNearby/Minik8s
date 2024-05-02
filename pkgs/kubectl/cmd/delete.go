package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/kubectl/api"
	"minik8s/utils"
	"os"
	"reflect"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <resource> <name>",
	Short: "Kubectl can delete resources and names",
	Long:  "Kubectl can delete resources and names",
	Run:   applyHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
	fileInfo, err := os.Stat(args[0])
	if err != nil {
		log.Fatal(err)
	}
	if fileInfo.IsDir() {
		log.Errorf("%s is not a file", args[0])
		cmd.Usage()
		return
	}
	fileContent, err := utils.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
		return
	}
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
	err = api.ParseApiObjectFromYamlFile(fileContent, object)
	if err != nil {
		log.Fatal(err)
	}
	err = controller.DeleteObject(objType, object.GetNameSpace(), object.GetNameSpace())

}
