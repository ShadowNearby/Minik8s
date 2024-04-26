package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubeclt"
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
	kind, err := kubeclt.GetApiKindFromYamlFile(fileContent)
	if err != nil {
		log.Fatal(err)
		return
	}
	structType, res := core.KindToStructType[kind]
	if !res {
		log.Error("Unsupported kind", kind)
		return
	}
	object := reflect.New(structType).Interface().(core.ApiObjectKind)
	err = kubeclt.ParseApiObjectFromYamlFile(fileContent, object)
	if err != nil {
		log.Fatal(err)
	}
	_url := utils.ParseUrlOne(kind, object.GetObjectName(), object.GetObjectNamespace())
	statusCode, err := utils.DelRequest(_url)
	if err != nil {
		log.Fatal(err)
	}
	if statusCode != 200 {
		log.Error("delete Pod", "create task request failed")
	} else {
		log.Info("delete Pod", "create task request success")
	}
}
