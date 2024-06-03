package serverless

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"net/http"
	"os"
	"testing"
)

func TestCreateFunction(t *testing.T) {
	url := fmt.Sprintf("http://%s:8090/api/v1/functions/", config.ClusterMasterIP)
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "function/function.json"))
	if err != nil {
		t.Errorf("read file error")
		return
	}
	var function core.Function
	err = utils.JsonUnMarshal(string(file), &function)
	if err != nil {
		t.Errorf("funtion.json format is wrong")
		return
	}
	code, info, err := utils.SendRequest("POST", url, file)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		return
	}
	if code != http.StatusOK {
		var infoType core.InfoType
		utils.JsonUnMarshal(info, &infoType)
		t.Errorf("internal error: %d: %s", code, infoType.Error)
	}

}

func TestTriggerFunction(t *testing.T) {
	url := fmt.Sprintf("http://%s:8090/api/v1/functions/%s/trigger", config.ClusterMasterIP, "serverless_app")
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "function/trigger.json"))
	if err != nil {
		t.Errorf("read file error")
		return
	}
	var trigger core.TriggerMessage
	_ = utils.JsonUnMarshal(string(file), &trigger)
	code, info, err := utils.SendRequest("POST", url, file)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		return
	}
	if code != http.StatusOK {
		var infoType core.InfoType
		utils.JsonUnMarshal(info, &infoType)
		t.Errorf("internal error: %d: %s", code, infoType.Error)
	}

}
