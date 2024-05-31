package kubeproxy

import (
	"minik8s/config"
	"minik8s/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Run() {

}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubeproxy",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&utils.CustomFormatter{})
		if err := config.InitConfig(cfgFile); err != nil {
			logrus.Fatalf("Error initializing config: %s", err.Error())
		}

		var serviceController ServiceController
		serviceController.Run()
		select {}
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
