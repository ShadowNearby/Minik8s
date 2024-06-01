package serverless

import (
	"bytes"
	"minik8s/pkgs/kubectl/cmd"
	"testing"
)

func TestFunctionRun(t *testing.T) {
	actual := new(bytes.Buffer)
	cmd.RootCommand.SetOut(actual)
	cmd.RootCommand.SetErr(actual)
	cmd.RootCommand.SetArgs([]string{"apply", "-f", "../../example/serverless/addfunc.yaml"})
	cmd.RootCommand.Execute()
	cmd.RootCommand.SetArgs([]string{"get", "functions"})
	cmd.RootCommand.Execute()
	cmd.RootCommand.SetArgs([]string{"trigger", "function", "-f", "../../example/eventTrigger.yaml"})
	cmd.RootCommand.Execute()
}
