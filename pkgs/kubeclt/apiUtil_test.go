package kubeclt

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
	"testing"
)

func TestGetApiKindFromYamlFile(t *testing.T) {
	content, err := utils.ReadFile("../testfile/yamlFile/service.yaml")
	if err != nil {
		t.Fatal(err)
	}
	// 把文件内容转换成API对象
	kind, err := GetApiKindFromYamlFile(content)

	if err != nil {
		t.Fatal(err)
	}
	if kind != "Service" {
		t.Fatal("kind should be Service")
	}
	t.Log(kind)
}
func TestGetObjectFromYamlFile(t *testing.T) {
	content, err := utils.ReadFile("../testfile/yamlFile/simple_deployment.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var pod core.Pod

	err = ParseApiObjectFromYamlFile(content, &pod)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pod)
}
