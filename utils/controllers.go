package utils

import core "minik8s/pkgs/apiobject"

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
