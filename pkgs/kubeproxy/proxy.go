package kubeproxy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func Init() {

}

func CreateService(ip string, port uint16) {
	//host := fmt.Sprintf("%s:%d", ip, port)

	bindCommandArg := fmt.Sprintf("addr add %s/24 dev flannel", ip)
	output, err := exec.Command("ip", bindCommandArg).CombinedOutput()
	if err != nil {
		log.Errorf("bind ip error: %s output: %s", err.Error(), output)
		return
	}
	natCommandArg := fmt.Sprintf("-t nat -A POSTROUTING -m ipvs --vaddr %s --vport %d -j MASQUERADE", ip, port)
	output, err = exec.Command("iptables", natCommandArg).CombinedOutput()
	if err != nil {
		log.Errorf("bind ip error: %s output: %s", err.Error(), output)
		return
	}

}
