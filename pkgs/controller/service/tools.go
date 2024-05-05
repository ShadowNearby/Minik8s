package service

import (
	"encoding/json"
	"errors"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/controller"
	"minik8s/pkgs/kubeproxy"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func FindUnusedIP() string {
	for i, used := range UsedIP {
		if i == 0 || used {
			continue
		}
		return fmt.Sprintf("%s%d", IPPrefix, i)
	}
	log.Errorf("No IP available")
	return ""
}

func FreeUsedIP(ip string) {
	indexs := strings.SplitN(ip, ",", -1)
	index := indexs[len(indexs)-1]
	ret, err := strconv.Atoi(index)
	if err != nil {
		log.Errorf("ip index to int error")
	}
	UsedIP[ret] = false
}

func MatchLabel(l map[string]string, r map[string]string) bool {
	for k, v := range l {
		if r[k] != v {
			return false
		}
	}
	return true
}

func FindDestPort(targetPort string, containers []core.Container) uint32 {
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.Name == targetPort {
				return p.ContainerPort
			}
		}
	}
	return 0
}

func PutSelector(service *core.Service) {
	key := fmt.Sprintf("%s-%s", service.MetaData.NameSpace, service.MetaData.Name)
	ServiceSelector[key] = &service.Spec.Selector
}

func DelSelector(service *core.Service) {
	key := fmt.Sprintf("%s-%s", service.MetaData.NameSpace, service.MetaData.Name)
	delete(ServiceSelector, key)
}

func GetSelector(service *core.Service) *core.Selector {
	key := fmt.Sprintf("%s-%s", service.MetaData.NameSpace, service.MetaData.Name)
	return ServiceSelector[key]
}

func CreateEndpointObject(service *core.Service) error {
	// get all pods
	response := controller.GetObject(core.ObjPod, service.MetaData.NameSpace, "")
	if response == "" {
		err := errors.New("cannot get pods")
		log.Errorf("get pod error: %s", err.Error())
		return err
	}
	pods := []core.Pod{}
	err := json.Unmarshal([]byte(response), &pods)
	if err != nil {
		log.Errorf("unmarshal pods error: %s", err.Error())
		return err
	}
	selectedPods := []core.Pod{}
	for _, pod := range pods {
		if MatchLabel(service.Spec.Selector.MatchLabels, pod.MetaData.Labels) {
			selectedPods = append(selectedPods, pod)
		}
	}

	for _, port := range service.Spec.Ports {
		endpoint := core.Endpoint{
			MetaData: core.MetaData{
				Name:      fmt.Sprintf("%s-%d", service.MetaData.Name, port.Port),
				NameSpace: service.MetaData.NameSpace,
			},
			Subsets: []core.EndpointSubset{
				{
					Addresses: []core.EndpointAddress{
						{
							IP: service.Spec.ClusterIP,
						},
					},
					Ports: []core.EndpointPort{
						{
							Port: port.Port,
						},
					},
				},
			},
		}
		err := controller.CreateObject(core.ObjEndPoint, endpoint.MetaData.Name, endpoint)
		if err != nil {
			log.Errorf("create endpoint error: %s", err.Error())
			return err
		}
		for _, pod := range selectedPods {
			destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
			kubeproxy.BindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
			log.Infof("create endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
		}
	}
	return nil
}

func DeleteEndpointObject(service *core.Service, pod *core.Pod) error {
	for _, port := range service.Spec.Ports {
		name := fmt.Sprintf("%s-%d", service.MetaData.Name, port.Port)
		namespace := service.MetaData.NameSpace
		if pod != nil {
			response := controller.GetObject(core.ObjEndPoint, namespace, name)
			if response == "" {
				err := errors.New("cannot get endpoint")
				log.Errorf("get endpoint error: %s", err.Error())
				return err
			}
			// endpoint := core.Endpoint{}
			// err := json.Unmarshal([]byte(response), &endpoint)
		}
		err := controller.DeleteObject(core.ObjEndPoint, namespace, name)
		if err != nil {
			log.Errorf("error in delete endpoint %s:%s", namespace, name)
			return err
		}
	}
	return nil
}

func GetAllServiceObject(namespace string) ([]core.Service, error) {
	response := controller.GetObject(core.ObjService, namespace, "")
	if response == "" {
		err := errors.New("cannot get services")
		log.Errorf("get services error: %s", err.Error())
		return nil, err
	}
	services := []core.Service{}
	err := json.Unmarshal([]byte(response), &services)
	if err != nil {
		log.Errorf("unmarshal services error: %s", err.Error())
		return nil, err
	}
	return services, nil
}
