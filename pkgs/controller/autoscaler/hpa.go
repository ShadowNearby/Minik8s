package autoscaler

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"math"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"time"
)

type HPAController struct {
}

func (h *HPAController) Apply(autoscaler core.HorizontalPodAutoscaler) error {
	//pods, err := findPods(autoscaler.Spec.ScaleTargetRef.Name, autoscaler.Spec.ScaleTargetRef.Namespace)
	pod, err := findPod(autoscaler.Spec.ScaleTargetRef.Name, autoscaler.Spec.ScaleTargetRef.Namespace)
	if err != nil {
		return err
	}
	updateHpa(&autoscaler, autoscaler.Spec.MinReplicas, 0) // todo, do not use pods, use rs
	return scaleUp(0, autoscaler.Spec.MinReplicas, pod, &autoscaler)
}

func (h *HPAController) Update(autoscaler core.HorizontalPodAutoscaler) error {
	pods, err := findPods(fmt.Sprintf("hpa-%s-", autoscaler.MetaData.Name), autoscaler.Spec.ScaleTargetRef.Namespace)
	if err != nil {
		return err
	}
	desiredReplicas := checkAndUpdateMetrics(pods, &autoscaler)
	if math.Abs(desiredReplicas-float64(len(pods))) <= 0.1 {
		desiredReplicas = float64(len(pods))
	} else {
		desiredReplicas = math.Ceil(desiredReplicas)
	}
	desiredInt := int(desiredReplicas)
	// update hpa
	updateHpa(&autoscaler, desiredInt, len(pods))
	if desiredInt > autoscaler.Spec.MaxReplicas {
		desiredInt = autoscaler.Spec.MaxReplicas
	} else if desiredInt < autoscaler.Spec.MinReplicas {
		desiredInt = autoscaler.Spec.MinReplicas
	}
	if len(pods) == desiredInt {
		return nil
	} else if len(pods) < desiredInt {
		return scaleUp(len(pods), desiredInt, pods[0], &autoscaler)
	} else {
		return scaleDown(len(pods), desiredInt, pods, &autoscaler)
	}
}

func updateHpa(autoscaler *core.HorizontalPodAutoscaler, desired, current int) {
	autoscaler.Status.DesiredReplicas = desired
	autoscaler.Status.CurrentReplicas = current
	autoscaler.Status.LastScaleTime = time.Now()
}

func findPod(name, namespace string) (core.Pod, error) {
	podTxt := utils.GetObject(core.ObjPod, namespace, name)
	var pod core.Pod
	err := utils.JsonUnMarshal(podTxt, &pod)
	return pod, err
}

// findPods find pods required by hpa
func findPods(prefix, namespace string) ([]core.Pod, error) { // TODO: pod name
	var pods []core.Pod
	err := storage.RangeGet(fmt.Sprintf("/pods/object/%s/%s", namespace, prefix), &pods)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// scaleDown should delete one or more pods
func scaleDown(currentReplica, desiredReplica int, currentPods []core.Pod, hpa *core.HorizontalPodAutoscaler) error {
	left := currentReplica - desiredReplica
	for i, pod := range currentPods {
		if i >= left {
			break
		}
		err := utils.DeleteObject(core.ObjPod, pod.GetNameSpace(), pod.MetaData.Name)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		hpa.Status.CurrentReplicas--
	}
	return nil
}

// scaleUp should create one or more pods
func scaleUp(currentReplica, desiredReplica int, template core.Pod, hpa *core.HorizontalPodAutoscaler) error {
	left := desiredReplica - currentReplica
	for i := 0; i < left; i++ {
		newTemplate := template
		newTemplate.MetaData.UUID = utils.GenerateUUID()
		newTemplate.MetaData.Name = fmt.Sprintf("hpa-%s-%s", hpa.MetaData.Name, utils.GenerateUUID())
		err := utils.CreateObject(core.ObjPod, newTemplate.GetNameSpace(), newTemplate)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		hpa.Status.CurrentReplicas++
	}
	return nil
}

// checkAndUpdateMetrics return desired replicas
func checkAndUpdateMetrics(pods []core.Pod, autoscaler *core.HorizontalPodAutoscaler) float64 {
	//desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
	var needCpu, needMem bool
	var desiredAvgCpu, desiredAvgMem int
	var desiredValCpu, desiredValMem uint64
	allMetrics := make([]core.Metrics, len(pods))
	for i, pod := range pods {
		allMetrics[i] = getPodMetrics(pod)
	}
	for _, resource := range autoscaler.Spec.Metrics.Resources {
		if resource.Name == "cpu" {
			needCpu = true
			desiredValCpu = resource.Target.Value
			desiredAvgCpu = resource.Target.AverageUtilization
		}
		if resource.Name == "memory" {
			needMem = true
			desiredValMem = resource.Target.Value
			desiredAvgMem = resource.Target.AverageUtilization
		}
	}
	var retVal float64 = 0
	if needCpu == true {
		var currentUtilization = 0
		var currentValue uint64 = 0
		for _, metric := range allMetrics {
			currentValue += metric.Resources[0].Target.Value
			currentUtilization += metric.Resources[0].Target.AverageUtilization
		}
		currentUtilization /= len(pods)
		currentValue /= uint64(len(pods))
		if desiredValCpu != 0 {
			retVal = float64(len(pods)) * (float64(currentValue) / float64(desiredValCpu))
		}
		if desiredAvgCpu != 0 {
			retVal = float64(len(pods)) * (float64(currentUtilization) / float64(desiredAvgCpu))
		}
		//return float64(len(pods)) * (float64(currentUtilization) / float64(desiredAvgCpu))
	}
	if needMem == true {
		var currentUtilization = 0
		var currentVal uint64 = 0
		for _, metric := range allMetrics {
			currentVal += metric.Resources[1].Target.Value
			currentUtilization += metric.Resources[1].Target.AverageUtilization
		}
		currentUtilization /= len(pods)
		currentVal /= uint64(len(pods))
		if desiredValMem != 0 {
			retVal = float64(len(pods)) * float64(currentVal) / float64(desiredValMem)
		}
		if desiredAvgMem != 0 {
			retVal = float64(len(pods)) * float64(currentUtilization) / float64(desiredAvgMem)
		}
		//return float64(len(pods)) * float64(currentUtilization) / float64(desiredAvgMem)
	}
	// update metrics
	autoscaler.Status.CurrentMetrics = allMetrics
	return retVal
}

func getPodMetrics(pod core.Pod) core.Metrics {
	if pod.Status.HostIP == "" {
		return core.Metrics{}
	}
	code, data, err := utils.SendRequest("GET",
		fmt.Sprintf("http://%s:10250/metrics/%s/%s",
			pod.Status.HostIP, pod.GetNameSpace(), pod.MetaData.Name),
		nil)
	if err != nil || code != 200 {
		return core.Metrics{}
	}
	var info core.InfoType
	utils.JsonUnMarshal(data, &info)
	var metrics core.Metrics
	utils.JsonUnMarshal(info.Data, &metrics)
	return metrics
}
