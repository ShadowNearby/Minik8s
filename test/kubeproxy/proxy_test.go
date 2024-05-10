package test

import (
	"fmt"
	"minik8s/pkgs/kubeproxy"
	"os/exec"
	"strconv"
	"testing"
)

func TestDocker(t *testing.T) {
	serverPorts := []uint32{20000, 20001, 20002}
	for i, serverPort := range serverPorts {
		createArgs := []string{"run", "-itd", "--name", fmt.Sprintf("hello-server%d", i), "-p",
			fmt.Sprintf("%d:%d", serverPort, serverPort),
			"hashicorp/http-echo", "-listen=:" + strconv.Itoa(int(serverPort)), fmt.Sprintf("-text=hello server%d", i)}
		output, err := exec.Command("docker", createArgs...).CombinedOutput()
		if err != nil {
			t.Fatalf("can not create image error: %s output: %s", err.Error(), output)
			return
		}

	}
	for i, serverPort := range serverPorts {
		curlArgs := []string{"curl", "-s", fmt.Sprintf("%s:%d", "localhost", serverPort)}
		output, err := exec.Command("curl", curlArgs...).CombinedOutput()
		if string(output) != fmt.Sprintf("hello server%d\n", i) {
			t.Errorf("output not match expect: %s, actual: %s", string(output), fmt.Sprintf("hello server%d\n", i))
		}
		if err != nil {
			t.Fatalf("can not curl error: %s output: %s", err.Error(), output)
		}
	}
	for i := range serverPorts {
		stopArgs := []string{"stop", fmt.Sprintf("hello-server%d", i)}
		rmArgs := []string{"rm", fmt.Sprintf("hello-server%d", i)}
		output, err := exec.Command("docker", stopArgs...).CombinedOutput()
		if err != nil {
			t.Fatalf("can not create image error: %s output: %s", err.Error(), output)
			return
		}
		output, err = exec.Command("docker", rmArgs...).CombinedOutput()
		if err != nil {
			t.Fatalf("can not create image error: %s output: %s", err.Error(), output)
			return
		}
	}
}

func TestIPVS(t *testing.T) {
	serverIP := "localhost"
	serverPorts := []uint32{20000, 20001, 20002}
	serviceIP := "10.10.0.1"
	servicePort := uint32(5678)
	err := kubeproxy.CreateService(serviceIP, servicePort)
	if err != nil {
		t.Fatalf("can not create service error: %s", err.Error())
		return
	}
	for i, serverPort := range serverPorts {
		createArgs := []string{"run", "-itd", "--name", fmt.Sprintf("hello-server%d", i), "-p",
			fmt.Sprintf("%d:%d", serverPort, serverPort),
			"hashicorp/http-echo", "-listen=:" + strconv.Itoa(int(serverPort)), fmt.Sprintf("-text=hello server%d", i)}
		output, err := exec.Command("docker", createArgs...).CombinedOutput()
		if err != nil {
			t.Fatalf("can not create image error: %s output: %s", err.Error(), output)
			return
		}
		err = kubeproxy.BindEndpoint(serviceIP, servicePort, serverIP, serverPort)
		if err != nil {
			t.Fatalf("can not create endpoint error: %s", err.Error())
			return
		}
	}
	curlArgs := []string{"curl", "-s", fmt.Sprintf("%s:%d", serviceIP, servicePort)}
	for i := len(serverPorts) - 1; i >= 0; i-- {
		output, err := exec.Command("curl", curlArgs...).CombinedOutput()
		if string(output) != fmt.Sprintf("hello server%d\n", i) {
			t.Errorf("output not match actual: %s, expect: %s", string(output), fmt.Sprintf("hello server%d\n", i))
		}
		if err != nil {
			t.Fatalf("can not curl error: %s output: %s", err.Error(), output)
		}
	}

	err = kubeproxy.DeleteService(serviceIP, servicePort)
	if err != nil {
		t.Fatalf("can not delete service error: %s", err.Error())
	}
	for i := range serverPorts {
		stopArgs := []string{"stop", fmt.Sprintf("hello-server%d", i)}
		rmArgs := []string{"rm", fmt.Sprintf("hello-server%d", i)}
		output, err := exec.Command("docker", stopArgs...).CombinedOutput()
		if err != nil {
			t.Fatalf("can not create image error: %s output: %s", err.Error(), output)
		}
		output, err = exec.Command("docker", rmArgs...).CombinedOutput()
		if err != nil {
			t.Fatalf("can not create image error: %s output: %s", err.Error(), output)
		}
	}
}
