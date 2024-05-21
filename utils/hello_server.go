package utils

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

func CreateHelloServer(port uint32, i int) error {
	createArgs := []string{"run", "-itd", "--name", fmt.Sprintf("hello-server%d", i), "-p",
		fmt.Sprintf("%d:%d", port, port),
		"hashicorp/http-echo:latest", "-listen=:" + strconv.Itoa(int(port)), fmt.Sprintf("-text=hello server%d", i)}
	output, err := exec.Command("docker", createArgs...).CombinedOutput()
	if err != nil {
		logrus.Errorf("can not create image error: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}

func TestHelloServer(addr string, i int) (bool, error) {
	curlArgs := []string{"-s", addr}
	output, err := exec.Command("curl", curlArgs...).CombinedOutput()
	if err != nil {
		logrus.Errorf("can not curl error: %s output: %s", err.Error(), output)
		return false, err
	}
	return true, nil
}

func DeleteHelloServer(port uint32, i int) error {
	stopArgs := []string{"stop", fmt.Sprintf("hello-server%d", i)}
	rmArgs := []string{"rm", fmt.Sprintf("hello-server%d", i)}
	output, err := exec.Command("docker", stopArgs...).CombinedOutput()
	if err != nil {
		logrus.Fatalf("can not create image error: %s output: %s", err.Error(), output)
		return err
	}
	output, err = exec.Command("docker", rmArgs...).CombinedOutput()
	if err != nil {
		logrus.Fatalf("can not create image error: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}
