package monitor

import (
	"minik8s/config"
	"minik8s/utils"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "monitor",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&utils.CustomFormatter{})
		logrus.SetLevel(logrus.InfoLevel)
		if err := config.InitConfig(cfgFile); err != nil {
			logrus.Fatalf("Error initializing config: %s", err.Error())
		}
		Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config/config.json", "config file (default is ./config/config.json)")
}
