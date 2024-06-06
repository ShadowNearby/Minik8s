package kubeproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"

	log "github.com/sirupsen/logrus"
)

func FindUnusedIP(namespace string, name string) string {
	url := fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/services/%s/clusterip", config.ClusterMasterIP, config.ApiServerPort, namespace, name)
	_, resp, err := utils.SendRequest("GET", url, []byte{})
	if err != nil || resp == "" {
		log.Errorf("error in get new ClusterIP")
		return ""
	}
	info := core.InfoType{}
	err = utils.JsonUnMarshal(resp, &info)
	if err != nil {
		log.Errorf("error in unmarshal clusterIP")
		return ""
	}
	if info.Data == "" {
		log.Errorf("error in get new clusterIP")
		return ""
	}
	return info.Data
}

func FreeUsedIP(namespace string, name string) {
	url := fmt.Sprintf("http://%s:%s/api/v1/namespaces/%s/services/%s/clusterip", config.ClusterMasterIP, config.ApiServerPort, namespace, name)
	_, _, err := utils.SendRequest("DELETE", url, []byte{})
	if err != nil {
		log.Errorf("error in free ClusterIP")
	}
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
		if utils.MatchLabel(service.Spec.Selector.MatchLabels, pod.MetaData.Labels) && pod.Status.Phase == core.PodPhaseRunning {
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
				BindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
				log.Infof("create endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
			}
			endpoint.Binds = append(endpoint.Binds, core.EndpointBind{
				ServicePort:  port.Port,
				Destinations: Destinations,
			})
		}
	} else if service.Spec.Type == core.ServiceTypeNodePort {
		NodeIP := constants.AllIP
		endpoint = core.Endpoint{
			MetaData: core.MetaData{
				Name:      service.MetaData.Name,
				Namespace: service.MetaData.Namespace,
			},
			ServiceClusterIP: NodeIP,
		}
		for _, port := range service.Spec.Ports {
			Destinations := []core.EndpointDestination{}
			for _, pod := range selectedPods {
				destPort := FindDestPort(port.TargetPort, pod.Spec.Containers)
				Destinations = append(Destinations, core.EndpointDestination{
					IP:   pod.Status.PodIP,
					Port: destPort,
				})
				BindEndpoint(NodeIP, port.NodePort, pod.Status.PodIP, destPort)
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
		BindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
		log.Infof("bind endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)

		for i, bind := range endpoint.Binds {
			if bind.ServicePort != port.Port {
				continue
			}
			endpoint.Binds[i].Destinations = append(bind.Destinations, core.EndpointDestination{
				IP:   pod.Status.PodIP,
				Port: destPort,
			})
		}
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
		UnbindEndpoint(service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)
		log.Infof("delete endpoint: %s:%d -> %s:%d", service.Spec.ClusterIP, port.Port, pod.Status.PodIP, destPort)

		dest := core.EndpointDestination{
			IP:   pod.Status.PodIP,
			Port: destPort,
		}
		newDestinations := []core.EndpointDestination{}
		for i, bind := range endpoint.Binds {
			if endpoint.Binds[i].ServicePort != port.Port {
				continue
			}
			for _, d := range bind.Destinations {
				if d != dest {
					newDestinations = append(newDestinations, dest)
				}
			}
			endpoint.Binds[i].Destinations = newDestinations
		}
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
