package cmd

import (
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubectl/api"
	"minik8s/utils"
	ctlutils "minik8s/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger <kind> -f <FILENAME>",
	Short: "Kubectl trigger command",
	Long:  "Kubectl trigger command, Usage: kubectl trigger <kind> (-f FILENAME)",
	Run:   trigger,
}

func trigger(cmd *cobra.Command, args []string) {
	log.Info("Trigger started")
	kind := strings.ToLower(args[0])
	if kind != "function" && kind != "workflow" {
		log.Errorln("invalid resource type, it should be function or workflow")
	}
	var objType core.ObjType
	if kind == "function" {
		objType = core.ObjFunction
	}
	if kind == "workflow" {
		objType = core.ObjWorkflow
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Error(err)
		return
	}
	if fileInfo.IsDir() {
		log.Errorf("%s is not a file", filePath)
		return
	}
	// 读取文件的内容
	fileContent, err := utils.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infoln("Object type: ", objType)
	name, _ := api.GetNameFromParamsFile(fileContent)
	paramsContent, err := api.GetParamsFromParamsFile(fileContent)
	if err != nil {

		log.Error(err)
	}
	log.Infoln("[Path]: ", filePath, "	[DATA]:	", paramsContent)
	info, _ := ctlutils.TriggerObject(objType, name, paramsContent)
	log.Infoln("the response: ", info)
}
