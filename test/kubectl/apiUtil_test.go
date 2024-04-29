package test

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubectl"
	"minik8s/utils"
	"os"
	"strings"
	"testing"
)

func TestGetApiKindFromYamlFile(t *testing.T) {
	rootpath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
		return
	}
	if !strings.HasSuffix(rootpath, "/test/kubectl") {
		t.Fatal("must in project root")
		return
	}
	content, err := utils.ReadFile(fmt.Sprintf("%s%s", rootpath, "/service.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	// 把文件内容转换成API对象
	kind, err := kubectl.GetApiKindFromYamlFile(content)

	if err != nil {
		t.Fatal(err)
	}
	if kind != "Service" {
		t.Fatal("kind should be Service")
	}
	t.Log(kind)
}
func TestGetObjectFromYamlFile(t *testing.T) {
	rootpath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
		return
	}
	if !strings.HasSuffix(rootpath, "/test/kubectl") {
		t.Fatal("file path error")
		return
	}
	content, err := utils.ReadFile(fmt.Sprintf("%s%s", rootpath, "/simple_deployment.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var pod core.Pod

	err = kubectl.ParseApiObjectFromYamlFile(content, &pod)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pod)
}
