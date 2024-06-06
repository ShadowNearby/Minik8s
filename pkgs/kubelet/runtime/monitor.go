package runtime

import (
	"errors"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/kubelet/resources"
	"minik8s/utils"

	v1 "github.com/google/cadvisor/info/v1"
	v2 "github.com/google/cadvisor/info/v2"
	logger "github.com/sirupsen/logrus"
)

func GetNodeState() (core.NodeMetrics, error) {
	nodeMetrics := core.NodeMetrics{}
	var machineAttr v2.Attributes
	var containerInfos v1.ContainerInfo
	if KubeletInstance.NumCores == 0 || KubeletInstance.MemCapacity == 0 {
		err := GetMachineAttr(&machineAttr)
		if err != nil {
			return nodeMetrics, err
		}
		KubeletInstance.NumCores = machineAttr.NumCores
		KubeletInstance.MemCapacity = machineAttr.MemoryCapacity
	} else {
		machineAttr.NumCores = KubeletInstance.NumCores
		machineAttr.MemoryCapacity = KubeletInstance.MemCapacity
	}
	err := getContainerInfos(&containerInfos)
	if err != nil {
		return nodeMetrics, err
	}
	numCores := machineAttr.NumCores
	deltaUsage := float64(containerInfos.Stats[1].Cpu.Usage.Total - containerInfos.Stats[0].Cpu.Usage.Total)
	deltaTime := float64(containerInfos.Stats[1].Timestamp.Sub(containerInfos.Stats[0].Timestamp))
	cpuUsage := deltaUsage / deltaTime
	cpuMetricVal := cpuUsage * float64(numCores) * 1000
	nodeMetrics.CPUUsage = cpuMetricVal
	memUsage := float64(containerInfos.Stats[0].Memory.Usage) / float64(machineAttr.MemoryCapacity)
	nodeMetrics.MemoryUsage = memUsage * 100.0
	fs := containerInfos.Stats[0].Filesystem
	var capacity, usage uint64
	for _, f := range fs {
		capacity += f.Limit
		usage += f.Usage
	}
	diskUsage := float64(usage) / float64(capacity)
	nodeMetrics.DiskUsage = diskUsage * 100.0
	logger.Info(nodeMetrics.CPUUsage, nodeMetrics.MemoryUsage, nodeMetrics.DiskUsage)
	return nodeMetrics, nil
}

func GetMachineAttr(attributes *v2.Attributes) error {
	path := "http://localhost:8080/api/v2.0/attributes"
	if code, info, _ := utils.SendRequest("GET", path, nil); code == 200 {
		utils.JsonUnMarshal(info, attributes)
		return nil
	} else {
		return errors.New(fmt.Sprintf("cannot successfully get machine attribute, code: %d", code))
	}
}

func getContainerInfos(infos *v1.ContainerInfo) error {
	path := "http://localhost:8080/api/v1.3/containers"
	if code, info, _ := utils.SendRequest("GET", path, nil); code == 200 {
		utils.JsonUnMarshal(info, infos)
		return nil
	} else {
		return errors.New(fmt.Sprintf("cannot successfully get container infos, code: %d", code))
	}
}

func GetPodMetrics(pod *core.Pod) (core.Metrics, error) {
	containerMetrics, err := resources.GetPodContainersMetrics(pod)
	if err != nil {
		logger.Error(err.Error())
		return core.Metrics{}, err
	}
	var allCpu, allMem float64
	for _, metric := range containerMetrics {
		allCpu += metric.CpuUsage
		allMem += metric.MemoryUsage
	}
	if KubeletInstance.NumCores == 0 || KubeletInstance.MemCapacity == 0 {
		var machineAttr v2.Attributes
		if KubeletInstance.NumCores == 0 || KubeletInstance.MemCapacity == 0 {
			GetMachineAttr(&machineAttr)
			KubeletInstance.NumCores = machineAttr.NumCores
			KubeletInstance.MemCapacity = machineAttr.MemoryCapacity
		}
	}
	var milliCore = uint64(allCpu) * 10 // milli-core
	var cpuUtilization = (int(milliCore) / (KubeletInstance.NumCores * 1000)) * 100
	var memUtilization = int((allMem / float64(KubeletInstance.MemCapacity)) * 100.0)
	resourceCpu := core.Resource{
		Name: "cpu",
		Target: core.ResourceTarget{
			Type:               "Utilization",
			Value:              milliCore,
			AverageUtilization: cpuUtilization,
		},
	}
	resourceMem := core.Resource{
		Name: "memory",
		Target: core.ResourceTarget{
			Type:               "Utilization",
			Value:              uint64(allMem),
			AverageUtilization: memUtilization,
		},
	}
	return core.Metrics{
		Type:      "Resource",
		Resources: []core.Resource{resourceCpu, resourceMem},
	}, nil
}
