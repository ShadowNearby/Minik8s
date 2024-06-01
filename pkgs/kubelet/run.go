package kubelet

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	kubeletcontroller "minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"

	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Run(config core.KubeletConfig, addr string) {
	runtime.KubeletInstance.InitKubelet(config)
	runtime.KubeletInstance.RegisterNode()
	for _, route := range kubeletcontroller.KubeletRouter {
		route.Register(runtime.KubeletInstance.Server)
	}
	go func() {
		for {
			err := runtime.KubeletInstance.Server.Run(addr)
			if err != nil {
				logger.Errorf("server run error: %s", err.Error())
			}
		}
	}()
	//go func() {
	//	for {
	//		for _, podConfig := range runtime.KubeletInstance.PodConfigMap {
	//			kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)
	//		}
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
	select {}
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubelet",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&utils.CustomFormatter{})
		var cfg = core.KubeletConfig{
			MasterIP:   config.ClusterMasterIP,
			MasterPort: config.ApiServerPort,
			Labels: map[string]string{
				"test": "haha",
				"app":  "nginx",
			},
		}
		Run(cfg, fmt.Sprintf("%s:%s", utils.GetIP(), config.NodePort))
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
