package main

import (
	log "github.com/sirupsen/logrus"
	"minik8s/pkgs/kubectl/cmd"
	test "minik8s/test/apiserver"
)

func main() {

	log.Println("Testing Kubectl")
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Error("Error executing commands: ", err)
	}
	//apiserver.ToolTest()
	//test.PodBasicTest()
	//test.PodLocalhostTest()
	//test.MetricsTest()
	//test.CreatePodTest()
	//test.InspectPod()
	test.ServerRun()
}
