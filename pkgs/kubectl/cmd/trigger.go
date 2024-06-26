package cmd

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubectl/api"
	"minik8s/utils"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger <resource> -f <file>",
	Short: "Kubectl can trigger functions or workflows",
	Long:  "Kubectl can trigger functions or workflows, usage trigger functions -f ./trigger.yaml",
	Run:   triggerHandler,
}

var resultCmd = &cobra.Command{
	Use:   "result <resource> <id>",
	Short: "Kubelet can get trigger result using id",
	Long:  "Kubelet can get trigger result using id, usage: result functions 1234567890",
	Run:   resultHandler,
}

func triggerHandler(cmd *cobra.Command, args []string) {
	var resourceKind string
	var resourceType core.ObjType
	logrus.Debugln(args)
	// check resource type
	if len(args) == 1 {
		resourceKind = strings.ToLower(args[0])
		for _, ty := range core.ObjTypeAll {
			if strings.Contains(ty, resourceKind) {
				resourceType = core.ObjType(ty)
				break
			}
		}
	} else {
		fmt.Printf("error: expect only one resource")
		return
	}
	if resourceType != core.ObjFunction && resourceType != core.ObjWorkflow && resourceType != core.ObjTask {
		fmt.Printf("error: only support functions type or workflows type")
		return
	}
	// parse file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}
	if fileInfo.IsDir() {
		logrus.Errorf("%s is not a file", filePath)
		return
	}
	fileContent, err := utils.ReadFile(filePath)
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}
	objType, err := api.GetObjTypeFromYamlFile(fileContent)
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}
	logrus.Infof("type: %s", objType)
	if core.ObjTask == objType {
		logrus.Infof("event trigger")
		taskMessage := core.PingSource{}
		err = yaml.Unmarshal(fileContent, &taskMessage)
		if err != nil {
			logrus.Fatal("cannot unmarshal file")
			return
		}
		url := fmt.Sprintf("http://%s:%s/api/v1/tasks", config.ClusterMasterIP, config.ApiServerPort)
		code, info, err := utils.SendRequest("POST", url, []byte(utils.JsonMarshal(taskMessage)))
		if err != nil {
			logrus.Errorf("send request failed; %s", err.Error())
			return
		}
		infoType := core.InfoType{}
		utils.JsonUnMarshal(info, &infoType)
		if code != http.StatusOK {
			fmt.Printf("code: %d, error: %s\n", code, infoType.Error)
			return
		} else {
			fmt.Printf("%s\n", infoType.Data)
			triggerResult := core.TriggerResult{
				ID:     taskMessage.ID,
				Result: infoType.Data,
			}
			utils.SaveTriggerResult(core.ObjFunction, triggerResult)
		}
		return
	}
	if core.ObjTrigger != objType {
		logrus.Fatal("expect trigger file, please set kind to trigger")
		return
	}
	triggerMsg := core.TriggerMessage{}
	err = yaml.Unmarshal(fileContent, &triggerMsg)
	if err != nil {
		logrus.Fatal("cannot unmarshal file")
		return
	}
	// send request
	var url string
	if resourceType == core.ObjFunction {
		url = fmt.Sprintf("http://%s:%s/api/v1/functions/%s/trigger", config.ClusterMasterIP, config.ApiServerPort, triggerMsg.Name)
	} else {
		url = fmt.Sprintf("http://%s:%s/api/v1/workflows/%s/trigger", config.ClusterMasterIP, config.ApiServerPort, triggerMsg.Name)
	}
	code, info, err := utils.SendRequest("POST", url, []byte(utils.JsonMarshal(triggerMsg)))
	if err != nil {
		logrus.Errorf("send request failed; %s", err.Error())
		return
	}
	infoType := core.InfoType{}
	utils.JsonUnMarshal(info, &infoType)
	if code != http.StatusOK {
		fmt.Printf("code: %d, error: %s\n", code, infoType.Error)
		return
	} else {
		fmt.Printf("%s\n", infoType.Data)
	}
}

func resultHandler(cmd *cobra.Command, args []string) {
	var resourceKind string
	var resourceType core.ObjType
	var resultId string
	logrus.Debugln(args)
	// check resource type
	if len(args) == 2 {
		resourceKind = strings.ToLower(args[0])
		for _, ty := range core.ObjTypeAll {
			if strings.Contains(ty, resourceKind) {
				resourceType = core.ObjType(ty)
				break
			}
		}
		resultId = args[1]
	} else {
		fmt.Printf("error: expect [resource] [id] format")
		return
	}
	if resourceType != core.ObjFunction && resourceType != core.ObjWorkflow && resourceType != core.ObjTask {
		fmt.Printf("error: only support functions type or workflows type")
		return
	}
	// send request to get result
	var url string
	if resourceType == core.ObjFunction || resourceType == core.ObjTask {
		url = fmt.Sprintf("http://%s:%s/api/v1/functions/result/%s", config.ClusterMasterIP, config.ApiServerPort, resultId)
	} else if resourceType == core.ObjWorkflow {
		url = fmt.Sprintf("http://%s:%s/api/v1/workflows/result/%s", config.ClusterMasterIP, config.ApiServerPort, resultId)
	}
	code, info, err := utils.SendRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("send request failed; %s", err.Error())
		return
	}
	infoType := core.InfoType{}
	utils.JsonUnMarshal(info, &infoType)
	if code != http.StatusOK {
		fmt.Printf("code: %d, error: %s\n", code, infoType.Error)
		return
	} else {
		fmt.Printf("%s\n", infoType.Data)
	}
}
