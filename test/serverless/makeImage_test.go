package serverless

import (
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/serverless/activator"
	"minik8s/pkgs/serverless/function"
	"minik8s/utils"
	"os/exec"
	"strings"
	"testing"
)

func TestCreateImage(t *testing.T) {
	activator.InitFunc("activate_example", utils.ExamplePath+"/serverless/func.py")
}
func TestRunImage(t *testing.T) {
	function.RunImage("activate_example")
}

func TestDeleteImage(t *testing.T) {
	err := function.DeleteImage("activate_example")
	if err != nil {
		t.Errorf("DeleteImage failed, error: %s", err)
	}
	// search the image
	cmd := exec.Command("docker", "images")
	out, err := cmd.Output()
	if err != nil {
		t.Errorf("DeleteImage failed, error: %s", err)
	}
	outputStr := string(out)
	imageName := fmt.Sprintf("%s:%s/%s:v1", config.LocalServerIp, config.ApiServerPort, "activate_example")
	if strings.Contains(outputStr, imageName) {
		t.Errorf("DeleteImage failed, image %s still exists", imageName)
	}
}
