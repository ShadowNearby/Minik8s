package test

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubectl/api"
	"minik8s/utils"
	"testing"
)

func TestGetApiKindFromYamlFile(t *testing.T) {
	content, err := utils.ReadFile(fmt.Sprintf("%s/replicaset.yaml", utils.ExamplePath))
	if err != nil {
		t.Fatal(err)
	}
	// 把文件内容转换成API对象
	kind, err := api.GetObjTypeFromYamlFile(content)
	if err != nil {
		t.Fatal(err)
	}
	if kind != core.ObjReplicaSet {
		t.Fatal("kind should be Service")
	}
	content, err = utils.ReadFile(fmt.Sprintf("%s/createPod.yaml", utils.ExamplePath))
	if err != nil {
		t.Fatal(err)
	}
	// 把文件内容转换成API对象
	kind, err = api.GetObjTypeFromYamlFile(content)
	if err != nil {
		t.Fatal(err)
	}
	if kind != core.ObjPod {
		t.Fatal("kind should be Pod")
	}
	content, err = utils.ReadFile(fmt.Sprintf("%s/service.yaml", utils.ExamplePath))
	// 把文件内容转换成API对象
	kind, err = api.GetObjTypeFromYamlFile(content)
	if err != nil {
		t.Fatal(err)
	}
	if kind != core.ObjService {
		t.Fatal("kind should be Pod")
	}
}
func TestGetObjectFromYamlFile(t *testing.T) {
	content, err := utils.ReadFile(fmt.Sprintf("%s/createPod.yaml", utils.ExamplePath))
	if err != nil {
		t.Fatal(err)
	}
	// 把文件内容转换成API对象
	kind, err := api.GetObjTypeFromYamlFile(content)
	if err != nil {
		t.Fatal(err)
	}
	obj := api.ParseApiObjectFromYamlFile(content, kind)
	log.Println(obj)
}
