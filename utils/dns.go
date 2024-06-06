package utils

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

func GenerateDNSPath(path string) string {
	splitPath := strings.Split(path, ".")
	result := config.DNSPathPrefix
	for i := len(splitPath) - 1; i >= 0; i-- {
		result = result + "/" + splitPath[i]
	}
	return result
}

func GetServiceAddress(namespace string, name string) (string, uint32, error) {
	if namespace == "" {
		namespace = "default"
	}
	content := GetObject(core.ObjService, namespace, name)
	service := &core.Service{}
	JsonUnMarshal(content, service)
	if len(service.Spec.Ports) < 1 {
		return "", 0, fmt.Errorf("service %s/%s has no port", namespace, name)
	}
	if service.Spec.ClusterIP != "" {
		return service.Spec.ClusterIP, service.Spec.Ports[0].Port, nil
	}
	NodeIP := GetIP()
	return NodeIP, service.Spec.Ports[0].NodePort, nil
}

func GenerateNginxFile(configs []core.DNSRecord) error {
	tmplpath := fmt.Sprintf("%s/%s", ConfigPath, "nginx.tmpl")
	confpath := fmt.Sprintf("%s/%s", ConfigPath, "nginx.conf")
	tmpl := template.Must(template.ParseFiles(tmplpath))
	Servers := []core.NginxServer{}
	for _, conf := range configs {
		locations := make([]core.NginxLocation, 0)
		for _, path := range conf.Paths {
			location := core.NginxLocation{
				Path: path.Path,
				IP:   path.IP,
				Port: path.Port,
			}
			locations = append(locations, location)
		}
		server := core.NginxServer{
			Addr:       config.NginxListenAddr,
			ServerName: conf.Host,
			Locations:  locations,
		}
		Servers = append(Servers, server)
	}
	config := core.NginxConf{
		Servers: Servers,
	}
	conffile, err := os.OpenFile(confpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Errorf("error in create file %s err: %s", confpath, err.Error())
		return err
	}
	err = tmpl.Execute(conffile, config)
	if err != nil {
		logrus.Errorf("error in exec template file err:%s", err.Error())
		return err
	}
	return nil
}

func StartNginx() error {
	confpath := fmt.Sprintf("%s/%s", ConfigPath, "nginx.conf")
	err := exec.Command("nginx", []string{"-c", confpath}...).Run()
	if err != nil {
		logrus.Errorf("error in start nginx %s", err.Error())
		return err
	}
	return nil
}

func StopNginx() error {
	err := exec.Command("nginx", []string{"-s", "stop"}...).Run()
	if err != nil {
		logrus.Errorf("error in stop nginx %s", err.Error())
		return err
	}
	return nil
}

func ReloadNginx() error {
	err := exec.Command("nginx", []string{"-s", "reload"}...).Run()
	if err != nil {
		logrus.Errorf("error in reload nginx %s", err.Error())
		return err
	}
	return nil
}

func UpdateNginx() error {
	response := GetObjectWONamespace(core.ObjDNS, "")
	dnsRecords := []core.DNSRecord{}
	err := JsonUnMarshal(response, &dnsRecords)
	if err != nil {
		logrus.Errorf("error in unmarshal dns records %s", err.Error())
		return err
	}

	err = GenerateNginxFile(dnsRecords)
	if err != nil {
		logrus.Errorf("error in generate nginx file %s", err.Error())
		return err
	}

	err = ReloadNginx()
	if err != nil {
		StopNginx()
		config.NginxStarted = false
		logrus.Errorf("error in reload nginx: %s", err.Error())
		return err
	}

	return nil
}

func AddCoreDns() error {
	originalData, err := os.ReadFile(config.TempResolvPath)
	if err != nil {
		logrus.Errorf("error in read resolv.conf file %s", err.Error())
		return err
	}

	newData := []byte(fmt.Sprintf("nameserver %s\n", config.NginxListenIP))
	newData = append(newData, originalData...)

	err = os.WriteFile(config.TempResolvPath, newData, 0644)
	if err != nil {
		logrus.Errorf("error in write resolv.conf file %s", err.Error())
		return err
	}
	return nil
}
