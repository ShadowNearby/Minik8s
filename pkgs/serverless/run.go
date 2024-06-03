package serverless

import (
	"minik8s/pkgs/serverless/autoscaler"
	"minik8s/utils"
)

func Run() {
	go autoscaler.PeriodicMetric(30)
	//go eventfilter.FunctionSync("functions")
	//go eventfilter.WorkFlowSync("workflowexecutors")
	utils.WaitForever()
}
