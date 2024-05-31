package main

import (
	"minik8s/pkgs/kubectl/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	cmd.Execute()
}
