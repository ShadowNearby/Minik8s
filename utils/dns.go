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

func GetServiceIP(namespace string, name string) string {
	if namespace == "" {
		namespace = "default"
	}
	content := GetObject(core.ObjService, namespace, name)
	service := &core.Service{}
	JsonUnMarshal(content, service)
	if service.Spec.ClusterIP != "" {
		return service.Spec.ClusterIP
	}
	NodeIP := GetIP()
	return NodeIP
}

func GenerateNginxFile(configs []core.DNSRecord) {
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
		logrus.Errorf("error in create file %s", confpath)
		return
	}
	err = tmpl.Execute(conffile, config)
	if err != nil {
		logrus.Errorf("error in exec template file")
	}
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
