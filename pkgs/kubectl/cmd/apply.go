package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "minik8s/pkgs/apiobject"

	"minik8s/pkgs/kubectl"
	"minik8s/utils"
	"os"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply can create apiObject in a declarative way",
	Long:  "Kubectl apply can create apiObject in a declarative way, usage kubectl apply -f [file]",
	Run:   applyHandler,
}

func applyHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
	// 检查参数是否是文件 读取文件
	fileInfo, err := os.Stat(args[0])
	if err != nil {
		log.Fatal(err)
		cmd.Usage()
		return
	}
	if fileInfo.IsDir() {
		log.Errorf("%s is not a file", args[0])
		cmd.Usage()
		return
	}
	// 读取文件的内容
	fileContent, err := utils.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	// 解析API对象的种类
	kind, err := kubectl.GetApiKindFromYamlFile(fileContent)

	if err != nil {
		log.Fatal(err)
		return
	}
	switch kind {
	case "Pod":
		log.Debug("[Case]:making apiObject apply pod")
		applyPodHandler(fileContent)
	case "Service":
		log.Debug("[Case]:making apiObject apply service")
		applyServiceHandler(fileContent)
	case "Job":
		log.Debug("[Case]:making apiObject apply job")
	case "Dns":
		log.Debug("[Case]:making apiObject apply dns")
	default:
		log.Debug("[Case]:unknown apiObject kind")
	}

}
func applyPodHandler(fileContent []byte) {
	var pod core.Pod
	err := kubectl.ParseApiObjectFromYamlFile(fileContent, &pod)
	log.Debug(pod)
	if err != nil {
		log.Error("apply Pod", "parse yaml failed", err.Error())
	}
	namespace := pod.MetaData.NameSpace
	_url := utils.ParseUrlMany("pod", namespace)
	statusCode, bodyJson, err := utils.PostRequestByTarget(_url, pod)
	if err != nil {
		log.Error("apply Pod", "post request failed", err.Error())
	}
	if statusCode != 200 {
		log.Error("apply Pod", "create task request failed")
	} else {
		log.Info("apply Pod", "create task request success", bodyJson)
		log.Infoln(pod)
	}
}
func applyServiceHandler(fileContent []byte) {
	var service core.Service
	err := kubectl.ParseApiObjectFromYamlFile(fileContent, &service)
	log.Debug(service)
	if err != nil {
		log.Error("apply Pod", "parse yaml failed", err.Error())
	}
	namespace := service.MetaData.NameSpace
	_url := utils.ParseUrlMany("service", namespace)
	statusCode, bodyJson, err := utils.PostRequestByTarget(_url, service)
	if err != nil {
		log.Error("apply Pod", "post request failed", err.Error())
	}
	if statusCode != 200 {
		log.Error("apply Pod", "create task request failed")
	} else {
		log.Info("apply Pod", "create task request success", bodyJson)
		log.Infoln(service)
	}

}
