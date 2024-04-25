package runtime

import (
	"github.com/gin-gonic/gin"
	core "minik8s/pkgs/apiobject"
	"time"
)

type Kubelet struct {
	PodMap   map[string]core.PodStatus
	IDtoName map[string]string
	Server   *gin.Engine
}

var KubeletInstance Kubelet

func (k *Kubelet) ContainerStart(podStatus *core.PodStatus, containerName string, containerID string) {
	//podStatus := k.PodMap[podConfig.MetaData.Name]
	ctStat := podStatus.Containers[containerID]
	ctStat.Ready = true
	ctStat.State.State = core.PhaseRunning
	ctStat.State.Started = time.Now()
	// write back
	podStatus.Containers[containerID] = ctStat
	// create IDtoName mapping
	if k.IDtoName == nil {
		k.IDtoName = make(map[string]string)
	}
	k.IDtoName[containerID] = containerName
	//k.PodMap[podConfig.MetaData.Name] = podStatus
}

func (k *Kubelet) InsertMap(key string, status core.PodStatus) {
	if k.PodMap == nil {
		k.PodMap = make(map[string]core.PodStatus)
	}
	k.PodMap[key] = status
}
