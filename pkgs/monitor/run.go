package monitor

import (
	"minik8s/config"
	"minik8s/utils"
	"time"

	"github.com/sirupsen/logrus"
)

func Run() {
	for {
		err := utils.GeneratePrometheusNodeFile()
		if err != nil {
			logrus.Errorf("error in generate node file %s", err.Error())
		}
		err = utils.GeneratePrometheusPodFile()
		if err != nil {
			logrus.Errorf("error in generate pod file %s", err.Error())
		}
		time.Sleep(config.PrometheusScrapeInterval)
	}
}
