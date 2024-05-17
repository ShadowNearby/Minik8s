package main

import (
	"minik8s/pkgs/apiserver"
	"minik8s/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})
	apiserver.Run()
}
