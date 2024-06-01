package apiserver

import (
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/controller/autoscaler"
	"minik8s/pkgs/controller/podcontroller"
	rsc "minik8s/pkgs/controller/replicaset"
	scheduler "minik8s/pkgs/controller/scheduler"
	"minik8s/pkgs/controller/service"
	"minik8s/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Run() {
	server := server.CreateAPIServer(config.DefaultEtcdEndpoints)
	var serviceController service.ServiceController
	go controller.StartController(&serviceController)
	var schedulerController scheduler.Scheduler
	go schedulerController.Run(constants.PolicyCPU)
	var podcontroller podcontroller.PodController
	go controller.StartController(&podcontroller)
	var replicaSet rsc.ReplicaSetController
	go controller.StartController(&replicaSet)
	var hpa autoscaler.HPAController
	go controller.StartController(&hpa)
	// start hpa background work
	go hpa.StartBackground()

	server.Run(fmt.Sprintf("%s:%s", config.ClusterMasterIP, config.ApiServerPort))
}

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&utils.CustomFormatter{})
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
