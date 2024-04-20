package cmd

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/testdata/imports/fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubeclt"
	"minik8s/utils"
	"os"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Kubectl apply can create apiObject in a declarative way",
	Long:  "Kubectl apply can create apiObject in a declarative way, usage kubectl apply -f [file]",
	Run:   applyHandler,
}

type ApplyOptions string

const (
	ApplyKindPod        ApplyOptions = "Pod"
	ApplyKindJob        ApplyOptions = "Job"
	ApplyKindService    ApplyOptions = "Service"
	ApplyKindReplicaset ApplyOptions = "Replicaset"
	ApplyKindDns        ApplyOptions = "Dns"
	ApplyKindHpa        ApplyOptions = "Hpa"
	ApplyKindFunc       ApplyOptions = "Function"
	ApplyKindWorkflow   ApplyOptions = "Workflow"
)

func applyHandler(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		return
	}
	// 检查参数是否是文件 读取文件
	fileInfo, err := os.Stat(args[0])
	if err != nil {
		fmt.Println(err.Error())
		cmd.Usage()
		return
	}
	if fileInfo.IsDir() {
		fmt.Println("file is a directory")
		cmd.Usage()
		return
	}
	// 读取文件的内容
	fileContent, err := utils.ReadFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 解析API对象的种类
	kind, err := kubeclt.GetApiKindFromYamlFile(fileContent)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch kind {
	case "Pod":
		log.Debug("[Case]:making apiObject create pod")
		var pod core.Pod
		kubeclt.ParseApiObjectFromYamlFile(fileContent, &pod)
		log.Debug(pod)

	case "Service":
		log.Debug("[Case]:making apiObject create service")
		var service core.Service
		kubeclt.ParseApiObjectFromYamlFile(fileContent, &service)
		log.Debug(service)
	default:
		log.Debug("[Case]:unknown apiObject kind")
	}

}
