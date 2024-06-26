package utils

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/containerd/containerd/namespaces"
	logger "github.com/sirupsen/logrus"
)

var nerdCtl, _ = exec.LookPath("nerdctl")
var ctrPath, _ = exec.LookPath("ctr")

type NerdCtl struct {
	namespace     string
	containerName string
	ctlType       string
}

const (
	NerdStop    string = "stop"
	NerdRm      string = "rm"
	NerdInspect string = "inspect"
	NerdCp      string = "cp"
)

type Ctr struct {
	ctrType       string
	ctrOp         string
	containerName string
	namespace     string
}

const (
	CtrSnapshot string = "snapshots"
)

const (
	CtrRm string = "rm"
)

func NerdRun(args ...string) (string, error) {
	res, err := exec.Command(nerdCtl, args...).CombinedOutput()
	return string(res), err
}

func NerdExec(ctl NerdCtl, args ...string) (string, error) {
	namespace := namespaces.Default
	if ctl.namespace != "" {
		namespace = ctl.namespace
	}
	if ctl.containerName == "" {
		return "", errors.New("expect container name")
	}
	containerName := ctl.containerName
	var cmd = make([]string, 0)
	cmd = append(cmd, "-n", namespace, ctl.ctlType)
	cmd = append(cmd, args...)
	cmd = append(cmd, containerName)
	logger.Debugf("exec: %s", cmd)
	res, err := exec.Command(nerdCtl, cmd...).CombinedOutput()
	return string(res), err
}

func NerdCopy(src string, dst string, namespace string) error {
	output, err := exec.Command(nerdCtl, NerdCp, src, dst, "--namespace", namespace).CombinedOutput()
	if err != nil {
		logger.Errorf("cp %s to %s failed: output: %s err: %s", src, dst, string(output), err.Error())
		return err
	}
	return nil
}

func CtrExec(ctr Ctr) (string, error) {
	namespace := namespaces.Default
	if ctr.namespace != "" {
		namespace = ctr.namespace
	}
	if ctr.containerName == "" {
		return "", errors.New("expect container name")
	}
	containerName := fmt.Sprintf("%s_%s", ctr.containerName, namespace)
	var cmd = make([]string, 0)
	cmd = append(cmd, ctr.ctrType, ctr.ctrOp, containerName)
	res, err := exec.Command(ctrPath, cmd...).CombinedOutput()
	return string(res), err
}
