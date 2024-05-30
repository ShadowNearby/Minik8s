package serverless

import (
	"minik8s/pkgs/serverless/activator"
	"minik8s/pkgs/serverless/function"
	"minik8s/utils"
	"os/exec"
	"strings"
	"testing"
)

//	func TestAll(t *testing.T) {
//		serverless.Run()sudo rm -rf /var/lib/docker
//
// sudo systemctl restart docker
//
//	}
func TestCreateImage(t *testing.T) {
	activator.InitFunc("activate", utils.ExamplePath+"/serverless/func.py")
}
func TestRunImage(t *testing.T) {
	function.RunImage("activate")
}

func TestDeleteImage(t *testing.T) {
	err := function.DeleteImage("activate")
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
	imageName := "localhost:5000/test:latest"
	if strings.Contains(outputStr, imageName) {
		t.Errorf("DeleteImage failed, image %s still exists", imageName)
	}
}
