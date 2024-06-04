package kubelet

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	kubeletcontroller "minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"minik8s/utils"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Run(kconfig core.KubeletConfig, addr string) {
	runtime.KubeletInstance.InitKubelet(kconfig)
	go func() {
		for {
			runtime.KubeletInstance.RegisterNode()
			time.Sleep(config.HeartbeatInterval)
		}
	}()
	for _, route := range kubeletcontroller.KubeletRouter {
		route.Register(runtime.KubeletInstance.Server)
	}
	go func() {
		for {
			err := runtime.KubeletInstance.Server.Run(addr)
			if err != nil {
				logrus.Errorf("server run error: %s", err.Error())
				return
			}
		}
	}()
	go func() {
		for {
			for i, podConfig := range runtime.KubeletInstance.PodConfigMap {
				if podConfig.Status.Condition != core.CondRunning {
					continue
				}
				err := kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)
				if err != nil {
					logrus.Errorf("error in inspect pod %s", podConfig.MetaData.Name)
				}
				if podConfig.Status.Condition != core.CondRunning {
					runtime.KubeletInstance.PodConfigMap[i] = podConfig
					logrus.Warnf("pod status changed")
					err := utils.SetObjectStatus(core.ObjPod, podConfig.MetaData.Namespace, podConfig.MetaData.Name, podConfig)
					if err != nil {
						logrus.Errorf("error in update pod status")
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
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
		if err := config.InitConfig(cfgFile); err != nil {
			logrus.Fatalf("Error initializing config: %s", err.Error())
		}
		host, _ := os.Hostname()
		var cfg = core.KubeletConfig{
			MasterIP:   config.ClusterMasterIP,
			MasterPort: config.ApiServerPort,
			Labels: map[string]string{
				"test": "haha",
				"app":  "nginx",
				"host": host,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config/config.json", "config file (default is ./config/config.json)")
}
