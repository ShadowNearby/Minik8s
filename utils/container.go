package utils

import (
	"context"
	"errors"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"time"

	v1 "github.com/containerd/cgroups/stats/v1"
	v2 "github.com/containerd/cgroups/v2/stats"
	"github.com/containerd/containerd"
	"github.com/containerd/typeurl/v2"
	logger "github.com/sirupsen/logrus"
)

//type ContainerWalker struct {
//	Client  containerd.Client
//	OnFound Found
//}
//
//type Found struct {
//	Container  containerd.Container
//	Req        string // The raw request string. name, short ID, or long ID.
//	MatchIndex int    // Begins with 0, up to MatchCount - 1.
//	MatchCount int    // 1 on exact match. > 1 on ambiguous match. Never be <= 0.
//}
//
//type OnFound func(ctx context.Context, found Found) error

func GenerateContainerSpec(pConfig core.Pod, cConfig core.Container, args ...string) core.ContainerdSpec {
	var cSpec = core.ContainerdSpec{
		Namespace:    pConfig.MetaData.Namespace,
		Image:        cConfig.Image,
		ID:           GenerateContainerIDByName(cConfig.Name, pConfig.MetaData.UUID),
		Name:         cConfig.Name,
		VolumeMounts: generateVolMountsMap(cConfig.VolumeMounts),
		Cmd:          cConfig.Cmd,
		Envs:         generateEnvList(cConfig.Env),
		Resource:     cConfig.Resources.Limit,
		Labels:       pConfig.MetaData.Labels,
		PullPolicy:   cConfig.ImagePullPolicy,
		PodName:      pConfig.MetaData.Name,
	}
	if len(args) > 0 && len(args[0]) > 0 {
		linuxNamespace := args[0]
		var (
			ns = map[string]string{
				"cgroup":  linuxNamespace + "cgroup",
				"uts":     linuxNamespace + "uts",
				"network": linuxNamespace + "net",
				"ipc":     linuxNamespace + "ipc",
			}
		)
		cSpec.LinuxNamespace = ns
	}
	return cSpec
}

func GenerateContainerLabel(podName string) map[string]string {
	return map[string]string{constants.MiniK8SPod: podName, constants.NerdctlName: podName}
}

// GetContainerStatus first return type of this function is nil or containerd.Status
func GetContainerStatus(container containerd.Container) (containerd.Status, error) {
	ctx := context.Background()
	task, err := container.Task(ctx, nil)
	if err != nil {
		logger.Errorf("cannot get task: %s", err.Error())
		return containerd.Status{
			Status:     containerd.Unknown,
			ExitStatus: 1,
			ExitTime:   time.Now(),
		}, err
	}
	status, err := task.Status(ctx)
	if err != nil {
		logger.Errorf("cannot get container status: %s", err.Error())
		return containerd.Status{
			Status:     containerd.Unknown,
			ExitStatus: 1,
			ExitTime:   time.Now(),
		}, err
	}
	return status, err
}

// GetContainerMetrics copy from containerd
func GetContainerMetrics(container containerd.Container) (core.ContainerMetrics, error) {
	ctx := context.Background()
	task, err := container.Task(ctx, nil)
	if err != nil {
		logger.Errorf("get running task error: %s", err.Error())
		return core.EmptyContainerMetrics, err
	}
	metric, err := task.Metrics(ctx)
	if err != nil {
		logger.Errorf("get task metrics error: %s", err.Error())
		return core.EmptyContainerMetrics, err
	}
	//var data interface{}
	var containerMetrics core.ContainerMetrics
	switch {
	case typeurl.Is(metric.Data, (*v2.Metrics)(nil)): // should be v2
		{
			data := &v2.Metrics{}
			if err := typeurl.UnmarshalTo(metric.Data, data); err != nil {
				return core.EmptyContainerMetrics, err
			}
			var ioMajor uint64 = 0
			for _, ioEntry := range data.Io.Usage {
				ioMajor += ioEntry.Major
			}
			containerMetrics.CpuUsage = data.CPU.UsageUsec
			containerMetrics.PidCount = data.Pids.Current
			containerMetrics.MemoryUsage = data.Memory.Usage
			containerMetrics.DiskUsage = ioMajor
			return containerMetrics, nil
		}
	case typeurl.Is(metric.Data, (*v1.Metrics)(nil)):
		{
			data := &v1.Metrics{}
			if err := typeurl.UnmarshalTo(metric.Data, data); err != nil {
				return core.EmptyContainerMetrics, err
			}
			containerMetrics.CpuUsage = data.CPU.Usage.Total
			containerMetrics.PidCount = data.Pids.Current
			containerMetrics.MemoryUsage = data.Memory.Usage.Usage
			return containerMetrics, nil
		}
	default:
		return core.EmptyContainerMetrics, errors.New("cannot convert metric data to cgroups.Metrics")
	}
	//marshaledJSON, err := json.MarshalIndent(data, "", "  ")
	//if err != nil {
	//	return err
	//}
	//logger.Infof("inspect data: %s", marshaledJSON) // 打印一下看看
	//return core.EmptyContainerMetrics, nil
}
