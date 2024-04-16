package utils

import (
	"errors"
	"fmt"
	"github.com/containerd/containerd/namespaces"
	"os/exec"
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
	res, err := exec.Command(nerdCtl, cmd...).CombinedOutput()
	return string(res), err
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
