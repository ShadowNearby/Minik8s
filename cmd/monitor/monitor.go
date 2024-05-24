package main

import (
	"minik8s/pkgs/monitor"
	"minik8s/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	monitor.Run()
}