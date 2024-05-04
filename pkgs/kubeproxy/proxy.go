package kubeproxy

import (
	"fmt"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func Init() {

}

// CreateService 创建虚拟节点，serviceIP 和 servicePort 代表虚拟节点的 IP 和端口，将虚拟节点添加到本地flannel网卡，添加SNAT功能
func CreateService(serviceIP string, servicePort uint32) error {
	createIPArgs := []string{"-A", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort), "-s", "rr"}
	output, err := exec.Command("ipvsadm", createIPArgs...).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to create ip: %s output: %s", err.Error(), output)
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
	deleteArgs := []string{"-D", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort)}
	output, err := exec.Command("ipvsadm", deleteArgs...).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to delete ip: %s output: %s", err.Error(), output)
		return err
	}
	unbindArgs := []string{"addr", "del", fmt.Sprintf("%s/24", serviceIP), "dev", "flannel.1"}
	output, err = exec.Command("ip", unbindArgs...).CombinedOutput()
	if err != nil {
		log.Errorf("unbind ip error: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}

// CreateEndpoint serviceIP 和 servicePort 代表虚拟节点的 IP 和端口，destIP 和 destPort 是真实节点的 IP 和端口，backend 是真实节点的 IP 地址。
// 在 IPVS 中添加服务的虚拟节点和真实节点的连接
func BindEndpoint(serviceIP string, servicePort uint32, destIP string, destPort uint32) error {
	addEndpointArgs := []string{"-a", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort), "-r", fmt.Sprintf("%s:%d", destIP, destPort), "-m"}
	output, err := exec.Command("ipvsadm", addEndpointArgs...).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to bind endpoint: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}

func UnbindEndpoint(serviceIP string, servicePort uint32, destIP string, destPort uint32) error {
	rmEndpointArgs := []string{"-d", "-t", fmt.Sprintf("%s:%d", serviceIP, servicePort), "-r", fmt.Sprintf("%s:%d", destIP, destPort)}
	output, err := exec.Command("ipvsadm", rmEndpointArgs...).CombinedOutput()
	if err != nil {
		log.Fatalf("failed to unbind endpoint: %s output: %s", err.Error(), output)
		return err
	}
	return nil
}
