package test

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"minik8s/pkgs/kubectl/cmd"
	"testing"
)

func TestApply(t *testing.T) {
	/* usable only when api-server is on */
	actual := new(bytes.Buffer)
	cmd.RootCommand.SetOut(actual)
	cmd.RootCommand.SetErr(actual)
	cmd.RootCommand.SetArgs([]string{"apply", "-f", "../example/createPod.yaml"})
	cmd.RootCommand.Execute()
	log.SetOutput(actual)
	expected := "This-is-command-a1"
	fmt.Print(actual.String())

	assert.Equal(t, actual.String(), expected, "actual is not expected")
}
