package main

import (
	"minik8s/pkgs/apiserver"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	apiserver.Run()
}
