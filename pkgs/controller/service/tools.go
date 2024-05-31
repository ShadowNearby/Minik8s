package service

import (
	"encoding/json"
	"errors"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/kubeproxy"
	"minik8s/utils"
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
	indexs := strings.SplitN(ip, ".", -1)
	index := indexs[len(indexs)-1]
	ret, err := strconv.Atoi(index)
	if err != nil {
		log.Errorf("ip index to int error ip %s index %s", ip, index)
	}
	UsedIP[ret] = false
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

func CreateEndpointObject(service *core.Service) error {
	// get all pods
	response := utils.GetObject(core.ObjPod, service.MetaData.Namespace, "")
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
		if utils.MatchLabel(service.Spec.Selector.MatchLabels, pod.MetaData.Labels) {
			selectedPods = append(selectedPods, pod)
		}
	}
	endpoint := core.Endpoint{}
	if service.Spec.Type == core.ServiceTypeClusterIP {
		endpoint = core.Endpoint{
			MetaData: core.MetaData{
				Name:      service.MetaData.Name,
				Namespace: service.MetaData.Namespace,
			},
			ServiceClusterIP: service.Spec.ClusterIP,
		}
		for _, port := range service.Spec.Ports {
			Destinations := []core.EndpointDestination{}
			for _, pod := range selectedPods {
				destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
				Destinations = append(Destinations, core.EndpointDestination{
					IP:   pod.Status.PodIP,
					Port: destPort,
				})
				kubeproxy.BindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
				log.Infof("create endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
			}
			endpoint.Binds = append(endpoint.Binds, core.EndpointBind{
				ServicePort:  port.Port,
				Destinations: Destinations,
			})
		}
	} else if service.Spec.Type == core.ServiceTypeNodePort {
		endpoint = core.Endpoint{
			MetaData: core.MetaData{
				Name:      service.MetaData.Name,
				Namespace: service.MetaData.Namespace,
			},
		}
		NodeIP := constants.Localhost
		for _, port := range service.Spec.Ports {
			Destinations := []core.EndpointDestination{}
			for _, pod := range selectedPods {
				destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
				Destinations = append(Destinations, core.EndpointDestination{
					IP:   NodeIP,
					Port: destPort,
				})
				kubeproxy.BindEndpoint(NodeIP, port.NodePort, pod.Status.PodIP, destPort)
				log.Infof("create endpoint: %s:%d -> %s:%d", NodeIP, port.NodePort, pod.Status.PodIP, destPort)
			}
			endpoint.Binds = append(endpoint.Binds, core.EndpointBind{
				ServicePort:  port.NodePort,
				Destinations: Destinations,
			})
		}
	}

	err = utils.CreateObject(core.ObjEndPoint, endpoint.MetaData.Namespace, endpoint)
	if err != nil {
		log.Errorf("create endpoint error: %s", err.Error())
		return err
	}
	return nil
}

func DeleteEndpointObject(service *core.Service) error {
	name := service.MetaData.Name
	namespace := service.MetaData.Namespace
	err := utils.DeleteObject(core.ObjEndPoint, namespace, name)
	if err != nil {
		log.Errorf("error in delete endpoint %s:%s", namespace, name)
		return err
	}
	return nil
}

func UpdateEndpointObjectByPodCreate(service *core.Service, pod *core.Pod) error {
	endpoint, err := GetEndpointObject(service)
	if err != nil {
		return err
	}
	if !utils.MatchLabel(service.Spec.Selector.MatchLabels, pod.MetaData.Labels) {
		return nil
	}
	for _, port := range service.Spec.Ports {
		destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
		kubeproxy.BindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
		log.Infof("bind endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)

		AddBinds(&endpoint.Binds, port.Port, core.EndpointDestination{
			IP:   pod.Status.PodIP,
			Port: destPort,
		})
	}
	err = utils.SetObject(core.ObjEndPoint, endpoint.MetaData.Namespace, endpoint.MetaData.Name, endpoint)
	if err != nil {
		log.Errorf("update endpoint error: %s", err.Error())
		return err
	}
	return nil
}

func UpdateEndpointObjectByPodDelete(service *core.Service, pod *core.Pod) error {
	endpoint, err := GetEndpointObject(service)
	if err != nil {
		return err
	}
	if !utils.MatchLabel(service.Spec.Selector.MatchLabels, pod.MetaData.Labels) {
		return nil
	}
	for _, port := range service.Spec.Ports {
		destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
		kubeproxy.UnbindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
		log.Infof("delete endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)

		RemoveBinds(&endpoint.Binds, port.Port, core.EndpointDestination{
			IP:   pod.Status.PodIP,
			Port: destPort,
		})
	}
	err = utils.SetObject(core.ObjEndPoint, endpoint.MetaData.Namespace, endpoint.MetaData.Name, endpoint)
	if err != nil {
		log.Errorf("update endpoint error: %s", err.Error())
		return err
	}
	return nil
}

func AddBinds(binds *[]core.EndpointBind, port uint32, dest core.EndpointDestination) {
	for _, bind := range *binds {
		if bind.ServicePort != port {
			continue
		}
		bind.Destinations = append(bind.Destinations, dest)
	}
}

func RemoveBinds(binds *[]core.EndpointBind, port uint32, dest core.EndpointDestination) {
	newDestinations := []core.EndpointDestination{}
	for _, bind := range *binds {
		if bind.ServicePort != port {
			continue
		}
		for _, d := range bind.Destinations {
			if d != dest {
				newDestinations = append(newDestinations, dest)
			}
		}
		bind.Destinations = newDestinations
	}
}

func GetEndpointObject(service *core.Service) (*core.Endpoint, error) {
	name := service.MetaData.Name
	namespace := service.MetaData.Namespace
	response := utils.GetObject(core.ObjEndPoint, namespace, name)
	if response == "" {
		err := errors.New("cannot get endpoint")
		log.Errorf("get endpoint error: %s", err.Error())
		return nil, err
	}
	endpoint := &core.Endpoint{}
	err := json.Unmarshal([]byte(response), &endpoint)
	if err != nil {
		log.Errorf("unmarshal endpoint error: %s", err.Error())
		return nil, err
	}
	return endpoint, nil
}

func GetAllServiceObject(namespace string) ([]core.Service, error) {
	response := utils.GetObject(core.ObjService, namespace, "")
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
