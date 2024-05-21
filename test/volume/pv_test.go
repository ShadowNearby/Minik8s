package test

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/volume"
	"minik8s/utils"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestVolumeBasic(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableQuote: true, ForceColors: true})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	pv := &core.PersistentVolume{}
	path := fmt.Sprintf("%s/%s", utils.ExamplePath, "pv.json")
	content, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("error in read file %s err: %s", path, err.Error())
		return
	}
	err = utils.JsonUnMarshal(string(content), pv)
	if err != nil {
		logrus.Errorf("error in unmarshal err: %s", err.Error())
		return
	}
	csiVolume, err := volume.CreateVolume(pv)
	if err != nil {
		logrus.Errorf("error in create volume err: %s", err.Error())
		return
	}
	err = volume.NodePublishVolume(csiVolume.VolumeId, config.CsiMntPath, pv)
	if err != nil {
		logrus.Errorf("error in publish volume err: %s", err.Error())
	}

	res := FsMountTestUtil(pv.Spec.Nfs.Share, config.CsiMntPath)
	if !res {
		logrus.Errorf("error in test monut: %s", err.Error())
	}

	err = volume.NodeUnpublishVolume(csiVolume.VolumeId, config.CsiMntPath)
	if err != nil {
		logrus.Errorf("error in unpublish volume err: %s", err.Error())
	}
	err = volume.DeleteVolume(csiVolume.VolumeId)
	if err != nil {
		logrus.Errorf("error in delete volume err: %s", err.Error())
	}
}

func TestMountUtil(t *testing.T) {
	res := FsMountTestUtil("/tmp", "/tmp")
	if !res {
		t.Fail()
	}
}

func FsMountTestUtil(src string, dst string) bool {
	content := "hello"
	srcPath := fmt.Sprintf("%s/%s", src, "hello.txt")
	srcFile, err := os.OpenFile(srcPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Errorf("error in open file %s err: %s", srcPath, err.Error())
		return false
	}
	defer srcFile.Close()
	_, err = srcFile.WriteString(content)
	if err != nil {
		logrus.Errorf("error in write file %s err: %s ", srcPath, err.Error())
		os.Remove(srcPath)
		return false
	}
	dstPath := fmt.Sprintf("%s/%s", dst, "hello.txt")
	raw, err := os.ReadFile(dstPath)
	if err != nil {
		logrus.Errorf("error in read file %s err: %s", dstPath, err.Error())
		os.Remove(srcPath)
		return false
	}
	if string(raw) == content {
		os.Remove(srcPath)
		return true
	}
	os.Remove(srcPath)
	return false
}
