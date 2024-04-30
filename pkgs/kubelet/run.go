package kubelet

import (
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/controller"
	"minik8s/pkgs/kubelet/runtime"
	"time"
)

var defaultAddr = ":10250"

func Run(config core.KubeletConfig, addr string) {
	runtime.KubeletInstance.InitKubelet(config)
	runtime.KubeletInstance.RegisterNode()
	go func() {
		for {
			err := runtime.KubeletInstance.Server.Run(addr)
			if err != nil {
				logger.Errorf("server run error: %s", err.Error())
			}
		}
	}()
	go func() {
		for {
			for _, podConfig := range runtime.KubeletInstance.PodConfigMap {
				kubeletcontroller.InspectPod(&podConfig, runtime.ExecProbe)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	select {}
}
