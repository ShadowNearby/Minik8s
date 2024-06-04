package scheduler

import (
	"math"
	core "minik8s/pkgs/apiobject"
)

type Scheduler struct{}

func (s *Scheduler) filterNode(nodes *[]core.Node, selector core.Selector) []core.Node {
	var res = make([]core.Node, 0)
	for _, node := range *nodes {
		for key, val := range selector.MatchLabels {
			if value, ok := node.MetaData.Labels[key]; ok != true || value != val {
				break
			}
			res = append(res, node)
		}
	}
	return res
}

func (s *Scheduler) scoreNode(nodes *[]core.Node) core.Node {
	scores := make([]uint64, len(*nodes))
	var lCPU, lMem, lPid, lDisk uint64
	var minimumIndex = 0
	var minimumScore uint64 = math.MaxUint64
	for _, node := range *nodes {
		if !node.Status.NetworkUnavailable || !node.Status.Ready {
			continue
		}
		if node.Status.CPUUsage > lCPU {
			lCPU = node.Status.CPUUsage
		}
		if node.Status.MemoryUsage > lMem {
			lMem = node.Status.MemoryUsage
		}
		if node.Status.PIDUsage > lPid {
			lPid = node.Status.PIDUsage
		}
		if node.Status.DiskUsage > lDisk {
			lDisk = node.Status.DiskUsage
		}
	}
	for i, node := range *nodes {
		if !node.Status.NetworkUnavailable || !node.Status.Ready {
			scores[i] = 0
			continue
		}
		scores[i] += lCPU - node.Status.CPUUsage
		scores[i] += lMem - node.Status.MemoryUsage
		scores[i] += lPid - node.Status.PIDUsage
		scores[i] += lDisk - node.Status.DiskUsage
		if scores[i] < minimumScore {
			minimumIndex = i
		}
	}
	return (*nodes)[minimumIndex]
}
