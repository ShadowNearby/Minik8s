package apiserver

import (
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/apiserver/heartbeat"
	"minik8s/pkgs/apiserver/server"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/controller/autoscaler"
	"minik8s/pkgs/controller/function"
	"minik8s/pkgs/controller/podcontroller"
	rsc "minik8s/pkgs/controller/replicaset"
	scheduler "minik8s/pkgs/controller/scheduler"
	"minik8s/utils"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Run() {
	server := server.CreateAPIServer(config.DefaultEtcdEndpoints)
	go server.Run(fmt.Sprintf("%s:%s", config.ClusterMasterIP, config.ApiServerPort))
	time.Sleep(3 * time.Second)
	var schedulerController scheduler.Scheduler
	go schedulerController.Run(constants.PolicyCPU)
	var podController podcontroller.PodController
	go controller.StartController(&podController)
	var replicaSet rsc.ReplicaSetController
	go controller.StartController(&replicaSet)
	go replicaSet.BackGroundTask()
	var hpa autoscaler.HPAController
	go controller.StartController(&hpa)
	// start hpa background work
	go hpa.StartBackground()
	var functionController function.FuncController
	go controller.StartController(&functionController)
	go functionController.ListenOtherChannels()
	var taskController function.TaskController
	go controller.StartController(&taskController)
	go taskController.StartTaskController()
	var workFlowController function.WorkFlowController
	go controller.StartController(&workFlowController)
	go workFlowController.StartController()
	// start heartbeat
	go heartbeat.Run()
	select {}
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/home/k8s/ly/minik8s/config/config.json", "config file (default is ./config/config.json)")
}
