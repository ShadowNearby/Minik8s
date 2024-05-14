package runtime

import (
	"errors"
	"fmt"
	v1 "github.com/google/cadvisor/info/v1"
	v2 "github.com/google/cadvisor/info/v2"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"
)

func GetNodeState() (core.NodeMetrics, error) {
	nodeMetrics := core.NodeMetrics{}
	var machineAttr v2.Attributes
	var containerInfos v1.ContainerInfo
	err := getMachineAttr(&machineAttr)
	if err != nil {
		return nodeMetrics, err
	}
	err = getContainerInfos(&containerInfos)
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

func getMachineAttr(attributes *v2.Attributes) error {
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
