package scheduler

import (
	"math"
	core "minik8s/pkgs/apiobject"
)

func roundRobinPolicy(idx int, candidates ...string) string {
	cLen := len(candidates)
	idx = idx % cLen
	return candidates[idx]
}

func cpuPolicy(candidates map[string]core.NodeMetrics) string {
	var bestCpuUsage = math.MaxFloat64
	var bestIp string
	for ip, metrics := range candidates {
		if metrics.CPUUsage < bestCpuUsage {
			bestIp = ip
		}
	}
	return bestIp
}

func memPolicy(candidates map[string]core.NodeMetrics) string {
	var bestMemUsage = math.MaxFloat64
	var bestIP string
	for ip, metrics := range candidates {
		if metrics.MemoryUsage < bestMemUsage {
			bestIP = ip
		}
	}
	return bestIP
}
func diskPolicy(candidates map[string]core.NodeMetrics) string {
	var bestDiskUsage = math.MaxFloat64
	var bestIP string
	for ip, metrics := range candidates {
		if metrics.DiskUsage < bestDiskUsage {
			bestIP = ip
		}
	}
	return bestIP
}
