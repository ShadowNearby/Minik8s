package test

import (
	"bytes"
	"minik8s/pkgs/kubectl/cmd"
	"minik8s/test/apiserver"
	"testing"
)

func TestApply(t *testing.T) {
	/* usable only when api-server is on */
	test.ServerRun()
	actual := new(bytes.Buffer)
	cmd.RootCommand.SetOut(actual)
	cmd.RootCommand.SetErr(actual)
	cmd.RootCommand.SetArgs([]string{"apply", "-f", "../files/createPod.yaml"})
	cmd.RootCommand.Execute()

	//expected := "This-is-command-a1"
	//fmt.Print(actual.String())
	//
	//assert.Equal(t, actual.String(), expected, "actual is not expected")
}
