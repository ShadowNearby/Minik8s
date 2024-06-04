package utils

import (
	"errors"
	"fmt"
	core "minik8s/pkgs/apiobject"
	"net/http"

	log "github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
)

func SendTriggerRequest(request core.TriggerRequest) (string, error) {
	code, info, err := SendRequest("POST", request.Url, request.Params)
	if err != nil {
		logger.Errorf("send reqeust error: %s", err.Error())
		return "", err
	}
	if code != http.StatusOK {
		var returnInfo core.InfoType
		JsonUnMarshal(info, returnInfo)
		logger.Errorf("request error: %s", returnInfo.Error)
		return "", errors.New(returnInfo.Error)
	}
	return info, nil
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

func ParseWorkStateMap(input map[string]core.WorkState) map[string]core.WorkState {
	result := make(map[string]core.WorkState)
	for key, workstate := range input {
		raw := JsonMarshal(workstate)
		var rawState core.RawState
		if err := JsonUnMarshal(raw, &rawState); err != nil {
			logger.Errorf("error parse workstate")
			return result
		}
		var state core.WorkState
		switch rawState.Type {
		case core.Task:
			var taskState core.TaskState
			if err := JsonUnMarshal(raw, &taskState); err != nil {
				log.Fatalf("Error unmarshaling TaskState: %v", err)
			}
			state = taskState
		case core.Fail:
			var failState core.FailState
			if err := JsonUnMarshal(raw, &failState); err != nil {
				log.Fatalf("Error unmarshaling FailState: %v", err)
			}
			state = failState
		case core.Choice:
			var choiceState core.ChoiceState
			if err := JsonUnMarshal(raw, &choiceState); err != nil {
				log.Fatalf("Error unmarshaling ChoiceState: %v", err)
			}
			state = choiceState
		default:
			log.Fatalf("Unknown state type: %v", rawState.Type)
		}

		result[key] = state
	}
	return result
}
