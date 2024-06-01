package autoscaler

/**
We implement horizontalPodAutoscaler in this file, it currently supports managing replicaset.
The HPA will inherent pods managed by replicaset and change their owner-reference when it's created, and when
HPA is updated, it will first check whether we can rescale and use the algorithm to calculate desired replicas
before scaling up or scaling down
By the way, it will start a background thread to periodically(5 minutes) update hpa to meet requirements
The rule of rescale is that the HPA didn't scale up/down in the last 5 minutes
The desired replicas refers to this equation: desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
Also notice that when abs(desiredReplicas - currentReplicas) < 0.1, we do not rescale
*/
import (
	"fmt"
	"math"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/utils"
	"time"

	logger "github.com/sirupsen/logrus"
)

type HPAController struct {
}

func (h *HPAController) StartBackground() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				h.BackgroundWork()
			}
		}
	}()
}

func (h *HPAController) GetChannel() string {
	return constants.ChannelHPA
}

func (h *HPAController) HandleCreate(message string) error {
	var hpa core.HorizontalPodAutoscaler
	utils.JsonUnMarshal(message, &hpa)
	return h.Apply(hpa)
}

func (h *HPAController) HandleUpdate(message string) error {
	hpas := make([]core.HorizontalPodAutoscaler, 2)
	utils.JsonUnMarshal(message, &hpas)
	oldHpa := hpas[0]
	hpa := hpas[1]
	if hpa.Spec.MinReplicas == oldHpa.Spec.MinReplicas &&
		hpa.Spec.MaxReplicas == oldHpa.Spec.MaxReplicas &&
		hpa.Spec.ScaleTargetRef == oldHpa.Spec.ScaleTargetRef {
		// there's no need to rescale
		logger.Errorf("no need to rescale")
		return nil
	}
	return h.Update(hpa)
}

func (h *HPAController) HandleDelete(message string) error {
	var hpa core.HorizontalPodAutoscaler
	utils.JsonUnMarshal(message, &hpa)
	return h.Delete(hpa)
}

func (h *HPAController) BackgroundWork() {
	// get all hpa
	var hpas []core.HorizontalPodAutoscaler
	hpasTxt := utils.GetObject(core.ObjHpa, "default", "")
	err := utils.JsonUnMarshal(hpasTxt, &hpas)
	if err != nil {
		return
	}
	// call update function
	for _, hpa := range hpas {
		err = h.Update(hpa)
		if err != nil {
			logger.Errorf("update hpa :%s error: %s", hpa.MetaData.Name, err.Error())
		}
	}
}

// Apply when create a hpa object, we use apply function
func (h *HPAController) Apply(autoscaler core.HorizontalPodAutoscaler) error {
	logger.Infof("[hpa-create] hpa name: %s", autoscaler.MetaData.Name)
	// find the replicaset
	var rs core.ReplicaSet
	rsTxt := utils.GetObject(core.ObjReplicaSet, "", autoscaler.Spec.ScaleTargetRef.Name)
	err := utils.JsonUnMarshal(rsTxt, &rs)
	if err != nil {
		logger.Errorf("unmarshal rs error: %s", err.Error())
		return err
	}
	// change replicaset owner_reference fo rs
	setRSController(&rs, autoscaler)
	// write rs information into storage
	err = utils.SetObjectStatus(core.ObjReplicaSet, "default", rs.MetaData.Name, rs)
	if err != nil {
		logger.Error("set rs status failed: ", err.Error())
		return err
	}
	// change pods managed by rs owner_reference to hpa
	pods, err := utils.FindRSPods(rs.MetaData.Name)
	for _, pod := range pods {
		setPodController(&pod, autoscaler)
		utils.SetObjectStatus(core.ObjPod, pod.MetaData.Namespace, pod.MetaData.Name, pod)
	}
	// decide whether to scale up or scale down
	desired := checkAndUpdateMetrics(pods, &autoscaler)
	desiredInt := getRealDesired(desired, len(pods), autoscaler)
	logger.Infof("desired replica number: %d, real number: %d", desiredInt, len(pods))
	if len(pods) == desiredInt {
		return nil
	} else if len(pods) < desiredInt {
		if len(pods) > 0 {
			err = scaleUp(len(pods), desiredInt, pods[0], &autoscaler)
		} else {
			pod := core.Pod{
				ApiVersion: autoscaler.ApiVersion,
				MetaData:   rs.Spec.Template.MetaData,
				Spec:       rs.Spec.Template.Spec,
				Status:     core.PodStatus{},
			}
			err = scaleUp(0, desiredInt, pod, &autoscaler)
		}
	} else {
		err = scaleDown(len(pods), desiredInt, pods, &autoscaler)
	}
	if err != nil {
		logger.Errorf("scale error: %s", err.Error())
		return err
	}
	// update last_scale time and other hpa information
	updateHpa(&autoscaler, desiredInt)
	// write hpa information into storage
	return utils.SetObjectStatus(core.ObjHpa, "default", autoscaler.MetaData.Name, autoscaler)
}

// Update when an hpa object updates, we use update
func (h *HPAController) Update(autoscaler core.HorizontalPodAutoscaler) error {
	logger.Infof("[hpa update] hpa name: %s", autoscaler.MetaData.Name)
	// check whether we can rescale
	if !canRescale(autoscaler.Status.LastScaleTime) {
		logger.Infof("interval too short, do not rescale")
		return nil
	}
	// get all pods managed by hpa
	pods, err := utils.FindHPAPods(autoscaler.MetaData.Name)
	if err != nil {
		logger.Errorf("get hpa pods error: %s", err.Error())
		return err
	}
	// decide whether to scale up or scale down
	desired := checkAndUpdateMetrics(pods, &autoscaler)
	desiredInt := getRealDesired(desired, len(pods), autoscaler)
	logger.Infof("desired replica number: %d, real number: %d", desiredInt, len(pods))
	if desiredInt == len(pods) {
		return nil
	} else if desiredInt < len(pods) {
		if len(pods) == 0 {
			// get replica first
			var rs core.ReplicaSet
			rsTxt := utils.GetObject(core.ObjReplicaSet, "", autoscaler.Spec.ScaleTargetRef.Name)
			utils.JsonUnMarshal(rsTxt, &rs)
			pod := core.Pod{
				ApiVersion: autoscaler.ApiVersion,
				MetaData:   rs.Spec.Template.MetaData,
				Spec:       rs.Spec.Template.Spec,
				Status:     core.PodStatus{},
			}
			err = scaleUp(len(pods), desiredInt, pod, &autoscaler)
		} else {
			err = scaleUp(len(pods), desiredInt, pods[0], &autoscaler)
		}
	} else {
		err = scaleDown(len(pods), desiredInt, pods, &autoscaler)
	}
	if err != nil {
		logger.Errorf("scale error: %s", err.Error())
		return err
	}
	// update last_scale time and other hpa information
	updateHpa(&autoscaler, desiredInt)
	// write hpa information into storage
	return utils.SetObjectStatus(core.ObjHpa, "default", autoscaler.MetaData.Name, autoscaler)
}

func (h *HPAController) Delete(autoscaler core.HorizontalPodAutoscaler) error {
	// get replicaset and delete
	utils.DeleteObject(core.ObjReplicaSet, "default", autoscaler.Spec.ScaleTargetRef.Name)
	// get pods and delete
	pods, err := utils.FindHPAPods(autoscaler.MetaData.Name)
	if err != nil {
		logger.Error("cannot find hpa pods: ", err.Error())
		return err
	}
	for _, pod := range pods {
		utils.DeleteObject(core.ObjPod, pod.MetaData.Namespace, pod.MetaData.Name)
	}
	return nil
}

// restrain the desired replica between [min, max] and mitigate the small difference between desired and actual replicas
func getRealDesired(desired float64, len int, autoscaler core.HorizontalPodAutoscaler) int {
	if math.Abs(desired-float64(len)) <= 0.1 {
		desired = float64(len)
	} else {
		desired = math.Ceil(desired)
	}
	desiredInt := int(desired)
	if desiredInt > autoscaler.Spec.MaxReplicas {
		desiredInt = autoscaler.Spec.MaxReplicas
	} else if desiredInt < autoscaler.Spec.MinReplicas {
		desiredInt = autoscaler.Spec.MinReplicas
	}
	return desiredInt
}

// check whether the time interval is enough for reschedule
func canRescale(lastScale time.Time) bool {
	duration := time.Now().Sub(lastScale)
	if duration < 5*time.Minute {
		return false
	}
	return true
}

// update desired and last scale time
func updateHpa(autoscaler *core.HorizontalPodAutoscaler, desired int) {
	autoscaler.Status.DesiredReplicas = desired
	autoscaler.Status.LastScaleTime = time.Now()
}

// set replicaset's owner-reference
func setRSController(rs *core.ReplicaSet, autoscaler core.HorizontalPodAutoscaler) {
	rs.MetaData.OwnerReference.Controller = true
	rs.MetaData.OwnerReference.ObjType = core.ObjHpa
	rs.MetaData.OwnerReference.Name = autoscaler.MetaData.Name
}

// set pod's owner-reference
func setPodController(pod *core.Pod, autoscaler core.HorizontalPodAutoscaler) {
	pod.MetaData.OwnerReference.Controller = true
	pod.MetaData.OwnerReference.ObjType = core.ObjHpa
	pod.MetaData.OwnerReference.Name = autoscaler.MetaData.Name
}

// scaleDown should delete one or more pods
func scaleDown(currentReplica, desiredReplica int, currentPods []core.Pod, hpa *core.HorizontalPodAutoscaler) error {
	left := currentReplica - desiredReplica
	for i, pod := range currentPods {
		if i >= left {
			break
		}
		err := utils.DeleteObject(core.ObjPod, pod.GetNamespace(), pod.MetaData.Name)
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
	setPodController(&template, *hpa)
	for i := 0; i < left; i++ {
		newTemplate := template
		newTemplate.MetaData.UUID = utils.GenerateUUID()
		newTemplate.MetaData.Namespace = "default"
		newTemplate.MetaData.Name = fmt.Sprintf("hpa-%s-%s", hpa.Spec.ScaleTargetRef.Name, utils.GenerateUUID(5))
		err := utils.CreateObject(core.ObjPod, newTemplate.GetNamespace(), newTemplate)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		hpa.Status.CurrentReplicas++
	}
	return nil
}

// checkAndUpdateMetrics get all pods' metrics and calculate desired replicas, also update HPA's status
func checkAndUpdateMetrics(pods []core.Pod, autoscaler *core.HorizontalPodAutoscaler) float64 {
	//desiredReplicas = ceil[currentReplicas * ( currentMetricValue / desiredMetricValue )]
	if len(pods) == 0 {
		// if there's no pod managed by hpa, we just return the minimum expect replicas
		return float64(autoscaler.Spec.MinReplicas)
	}
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
	var metricsNum = 0
	// we use the average desired pod number as return value
	if needCpu == true {
		var currentUtilization = 0
		var currentValue uint64 = 0
		var skipped = 0
		for _, metric := range allMetrics {
			if len(metric.Resources) < 1 {
				skipped += 1
			} else {
				currentValue += metric.Resources[0].Target.Value
				currentUtilization += metric.Resources[0].Target.AverageUtilization
			}
		}
		if len(pods) == skipped {
			logger.Error("not cpu information")
		} else {
			length := len(pods) - skipped
			currentUtilization /= length
			currentValue /= uint64(length)
			if desiredValCpu != 0 {
				retVal += float64(length) * (float64(currentValue) / float64(desiredValCpu))
				metricsNum++
			}
			if desiredAvgCpu != 0 {
				retVal += float64(length) * (float64(currentUtilization) / float64(desiredAvgCpu))
				metricsNum++
			}
		}
		//return float64(len(pods)) * (float64(currentUtilization) / float64(desiredAvgCpu))
	}
	if needMem == true {
		var currentUtilization = 0
		var currentVal uint64 = 0
		var skipped = 0
		for _, metric := range allMetrics {
			if len(metric.Resources) < 2 {
				skipped += 1
			} else {
				currentVal += metric.Resources[1].Target.Value
				currentUtilization += metric.Resources[1].Target.AverageUtilization
			}
		}
		length := len(pods) - skipped
		currentUtilization /= length
		currentVal /= uint64(length)
		if desiredValMem != 0 {
			retVal += float64(length) * float64(currentVal) / float64(desiredValMem)
			metricsNum++
		}
		if desiredAvgMem != 0 {
			retVal += float64(length) * float64(currentUtilization) / float64(desiredAvgMem)
			metricsNum++
		}
		//return float64(len(pods)) * float64(currentUtilization) / float64(desiredAvgMem)
	}
	// update metrics
	if metricsNum == 0 {
		return float64(autoscaler.Spec.MinReplicas)
	}
	autoscaler.Status.CurrentMetrics = allMetrics
	retVal = retVal / float64(metricsNum)
	return retVal
}

func getPodMetrics(pod core.Pod) core.Metrics {
	if pod.Status.HostIP == "" {
		return core.Metrics{}
	}
	code, data, err := utils.SendRequest("GET",
		fmt.Sprintf("http://%s:10250/metrics/%s/%s",
			pod.Status.HostIP, pod.GetNamespace(), pod.MetaData.Name),
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
