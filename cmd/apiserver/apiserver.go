package main

import (
	"minik8s/pkgs/apiserver"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logrus.SetReportCaller(true)
	apiserver.Run()
}
