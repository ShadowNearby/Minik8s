package test

import (
	"bytes"
	"minik8s/pkgs/kubectl/cmd"
	"testing"
)

func TestGet(t *testing.T) {
	/* usable only when api-server is on */
	actual := new(bytes.Buffer)
	cmd.RootCommand.SetOut(actual)
	cmd.RootCommand.SetErr(actual)
	cmd.RootCommand.SetArgs([]string{"get", "pods"})
	cmd.RootCommand.Execute()

	//expected := "This-is-command-a1"
	//fmt.Print(actual.String())
	//
	//assert.Equal(t, actual.String(), expected, "actual is not expected")
}
