package utils

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"net/http"
)

func SendTriggerRequest(request core.TriggerRequest) error {
	code, info, err := SendRequest("POST", request.Url, request.Params)
	if err != nil {
		logger.Errorf("send reqeust error: %s", err.Error())
		return err
	}
	if code != http.StatusOK {
		var returnInfo core.InfoType
		JsonUnMarshal(info, returnInfo)
		logger.Errorf("request error: %s", returnInfo.Error)
		return errors.New(returnInfo.Error)
	}
	return nil
}

// GenerateRSConfig  Generate Replicaset to save function*/
func GenerateRSConfig(name string, namespace string, image string, replicas int) *core.ReplicaSet {
	return &core.ReplicaSet{
		Kind:       string(core.ObjReplicaSet),
		ApiVersion: "v1",
		MetaData: core.MetaData{
			Name:      name,
			Namespace: namespace,
			OwnerReference: core.OwnerReference{
				Name:       name,
				ObjType:    core.ObjFunction,
				Controller: true,
			},
		},
		Spec: core.ReplicaSetSpec{
			Replicas: replicas,
			Selector: core.Selector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: core.ReplicaSetTemplate{
				MetaData: core.MetaData{
					Name:      name,
					Namespace: namespace,
					Labels:    map[string]string{"app": name},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            name,
							Image:           fmt.Sprintf("docker.io/%s", image),
							ImagePullPolicy: core.PullIfNeeds,
							Ports: []core.PortConfig{
								{
									ContainerPort: 80,
									Protocol:      "TCP",
									Name:          "p1",
								},
							},
							Cmd: []string{
								"python3",
								"server.py",
							},
						},
					},
				},
			},
		},
		Status: core.ReplicaSetStatus{
			RealReplicas: 0,
		},
	}
}
