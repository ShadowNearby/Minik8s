package utils

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
)

// FilterOwner give original pods, returns pods owned by controller<kind-namespace-name>
func FilterOwner(origin *[]core.Pod, name string, kind core.ObjType) []core.Pod {
	result := make([]core.Pod, 0)
	for _, pod := range *origin {
		or := pod.MetaData.OwnerReference
		if or.Controller &&
			or.ObjType == kind &&
			or.Name == name {
			result = append(result, pod)
		}
	}
	return result
}

func FindRSPods(rsName string, namespace ...string) ([]core.Pod, error) {
	// rsNamespace should be default
	// get all pods
	var pods []core.Pod
	ns := ""
	if len(namespace) > 0 {
		ns = namespace[0]
	}
	podsTxt := GetObject(core.ObjPod, ns, "")
	if podsTxt == "" {
		logger.Debugf("not pods found")
		return nil, nil
	}
	JsonUnMarshal(podsTxt, &pods)
	// filter pods with this rs owner-reference
	return FilterOwner(&pods, rsName, core.ObjReplicaSet), nil
}

func FindHPAPods(hpaName string) ([]core.Pod, error) {
	var pods []core.Pod
	podsTxt := GetObject(core.ObjPod, "", "")
	if podsTxt == "" {
		logger.Errorf("cannot find hpa pods")
		return nil, nil
	}
	JsonUnMarshal(podsTxt, &pods)
	return FilterOwner(&pods, hpaName, core.ObjHpa), nil
}

func FindFunctionRs(funcName string) (core.ReplicaSet, error) {
	var rs core.ReplicaSet
	rsTxt := GetObject(core.ObjReplicaSet, "default", funcName)
	JsonUnMarshal(rsTxt, &rs)
	// check the owner-reference
	or := rs.MetaData.OwnerReference
	if or.Controller != true || or.ObjType != core.ObjFunction || or.Name != funcName {
		return core.ReplicaSet{}, errors.New("the owner reference is wrong")
	}
	return rs, nil
}
