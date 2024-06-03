package utils

import (
	"context"
	"errors"
	"fmt"
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
	CheckContainerMataData(&cConfig)
	var cSpec = core.ContainerdSpec{
		Namespace:    pConfig.MetaData.Namespace,
		Image:        cConfig.Image,
		ID:           GenerateContainerIDByName(cConfig.Name, pConfig.MetaData.UUID),
		Name:         GenerateContainerName(pConfig.MetaData.Name, cConfig.Name),
		VolumeMounts: generateVolMountsMap(cConfig.VolumeMounts),
		Cmd:          cConfig.Cmd,
		Args:         cConfig.Args,
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

func GenerateContainerName(podName string, cName string) string {
	return fmt.Sprintf("%s-%s", podName, cName)
}

func GenerateContainerLabel(podName, containerName, namespace string) map[string]string {
	return map[string]string{
		constants.MiniK8SPod:       podName,
		constants.NerdctlName:      containerName,
		constants.MiniK8SNamespace: namespace,
	}
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
	var containerMetrics core.ContainerMetrics
	var cpuUsageUsec, memoryUsageBytes uint64
	var initialTime, finalTime time.Time
	for i := 0; i < 2; i++ {
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
		if i == 0 {
			initialTime = time.Now()
		} else {
			finalTime = time.Now()
		}
		switch {
		case typeurl.Is(metric.Data, (*v2.Metrics)(nil)): // should be v2
			{
				data := &v2.Metrics{}
				if err := typeurl.UnmarshalTo(metric.Data, data); err != nil {
					return core.EmptyContainerMetrics, err
				}
				if i == 0 {
					cpuUsageUsec = data.CPU.UsageUsec
					memoryUsageBytes = data.Memory.Usage
				} else {
					cpuUsageUsec = data.CPU.UsageUsec - cpuUsageUsec
					memoryUsageBytes = (data.Memory.Usage + memoryUsageBytes) / 2
				}
			}
		case typeurl.Is(metric.Data, (*v1.Metrics)(nil)):
			{
				data := &v1.Metrics{}
				if err := typeurl.UnmarshalTo(metric.Data, data); err != nil {
					return core.EmptyContainerMetrics, err
				}
				if i == 0 {
					cpuUsageUsec = data.CPU.Usage.Total
					memoryUsageBytes = data.Memory.Usage.Usage
				} else {
					cpuUsageUsec = data.CPU.Usage.Total - cpuUsageUsec
					memoryUsageBytes = (data.Memory.Usage.Usage + memoryUsageBytes) / 2
				}
			}
		default:
			return core.EmptyContainerMetrics, errors.New("cannot convert metric data to cgroups.Metrics")
		}
		if i == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
	elapsedTime := finalTime.Sub(initialTime).Microseconds()

	// 计算CPU,MEM使用率
	cpuUsageRate := (float64(cpuUsageUsec) / float64(elapsedTime)) * 100
	memUsageVal := float64(memoryUsageBytes)
	containerMetrics.CpuUsage = cpuUsageRate
	containerMetrics.MemoryUsage = memUsageVal
	return containerMetrics, nil
}

func CheckContainerMataData(container *core.Container) {
	if container.Name == "" {
		container.Name = GenerateUUID()
	}
}
