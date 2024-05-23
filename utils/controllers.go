package utils

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
)

// FilterOwner give original pods, returns pods owned by controller<kind-namespace-name>
func FilterOwner(origin *[]core.Pod, namespace, name string, kind core.ObjType) []core.Pod {
	result := make([]core.Pod, 0)
	for _, pod := range *origin {
		or := pod.MetaData.OwnerReference
		if or.Controller == true &&
			or.ObjType == kind &&
			or.Name == name {
			result = append(result, pod)
		}
	}
	return result
}

func FindRSPods(rsName string) ([]core.Pod, error) {
	// rsNamespace should be default
	// get all pods
	var pods []core.Pod
	podsTxt := GetObject(core.ObjPod, "", "")
	if podsTxt == "" {
		logger.Debugf("not pods found")
		return nil, nil
	}
	JsonUnMarshal(podsTxt, &pods)
	// filter pods with this rs owner-reference
	return FilterOwner(&pods, "default", rsName, core.ObjReplicaSet), nil
}

func FindHPAPods(hpaName string) ([]core.Pod, error) {
	var pods []core.Pod
	podsTxt := GetObject(core.ObjPod, "", "")
	if podsTxt == "" {
		return nil, errors.New("cannot get pods")
	}
	JsonUnMarshal(podsTxt, &pods)
	return FilterOwner(&pods, "default", hpaName, core.ObjHpa), nil
}
