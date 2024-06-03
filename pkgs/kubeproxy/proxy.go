package kubeproxy

import (
	"fmt"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// CreateService 创建虚拟节点，serviceIP 和 servicePort 代表虚拟节点的 IP 和端口，将虚拟节点添加到本地flannel网卡，添加SNAT功能
func CreateService(serviceIP string, servicePort uint32) error {
	log.Infof("create service on %s:%d", serviceIP, servicePort)
	createIPArgs := []string{"-A", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort), "-s", "rr"}
	output, err := exec.Command("ipvsadm", createIPArgs...).CombinedOutput()
	if err != nil {
		log.Errorf("failed to create ip: %s output: %s\naddr:%s:%d", err.Error(), output, serviceIP, servicePort)
		return err
	}

	bindCommandArg := []string{"addr", "add", fmt.Sprintf("%s/24", serviceIP), "dev", "flannel.1"}
	output, err = exec.Command("ip", bindCommandArg...).CombinedOutput()
	if err != nil {
		log.Errorf("bind ip error: %s output: %s", err.Error(), output)
		return err
	}
	natCommandArg := []string{"-t", "nat", "-A", "POSTROUTING", "-m", "ipvs", "--vaddr", serviceIP, "--vport", strconv.Itoa(int(servicePort)), "-j", "MASQUERADE"}
	output, err = exec.Command("iptables", natCommandArg...).CombinedOutput()
	if err != nil {
		log.Errorf("add endpoint error: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}

func DeleteService(serviceIP string, servicePort uint32) error {
	log.Infof("delete service on %s:%d", serviceIP, servicePort)
	deleteArgs := []string{"-D", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort)}
	output, err := exec.Command("ipvsadm", deleteArgs...).CombinedOutput()
	if err != nil {
		log.Errorf("failed to delete ip: %s output: %s serviceIP: %s servicePort: %d", err.Error(), output, serviceIP, servicePort)
		return err
	}
	unbindArgs := []string{"addr", "del", fmt.Sprintf("%s/24", serviceIP), "dev", "flannel.1"}
	output, err = exec.Command("ip", unbindArgs...).CombinedOutput()
	if err != nil {
		log.Errorf("unbind ip error: %s serviceIP: %s output: %s", err.Error(), serviceIP, output)
		return err
	}
	return nil
}

// CreateEndpoint serviceIP and servicePort represent the virtual node's IP and port,
// while destIP and destPort are the real node's IP and port. The backend is the IP address of the real node.
// Add the connection between the virtual node and the real node in IPVS for the service.
func BindEndpoint(serviceIP string, servicePort uint32, destIP string, destPort uint32) error {
	log.Infof("bind serviceIP %s servicePort %d destIP %s destPort %d", serviceIP, servicePort, destIP, destPort)
	addEndpointArgs := []string{"-a", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort), "-r", fmt.Sprintf("%s:%d", destIP, destPort), "-m"}
	output, err := exec.Command("ipvsadm", addEndpointArgs...).CombinedOutput()
	if err != nil {
		log.Errorf("failed to bind endpoint: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}

func UnbindEndpoint(serviceIP string, servicePort uint32, destIP string, destPort uint32) error {
	rmEndpointArgs := []string{"-d", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort), "-r", fmt.Sprintf("%s:%d", destIP, destPort)}
	output, err := exec.Command("ipvsadm", rmEndpointArgs...).CombinedOutput()
	if err != nil {
		log.Errorf("failed to unbind endpoint: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}
