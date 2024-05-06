package core

/*
	an basic example of a autoscaler apiobject:

apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:

	name: hpa-practice

spec:

	minReplicas: 3  # 最小pod数量
	maxReplicas: 6  # 最大pod数量
	metrics:
	- pods:
	    metric:
	      name: k8s_pod_rate_cpu_core_used_limit
	    target:
	      averageValue: "80"
	      type: AverageValue
	  type: Pods
	workload:   # 指定要控制的deploy
	  apiVersion:  apps/v1
	  kind: Pod
	  name: deploy-practice
	behavior:
	  scaleDown:
	    policies:
	    - type: Percent
	      value: 10
	      periodSeconds: 60 # 每分钟最多10%
*/
type HorizontalPodAutoscale struct {
	APIVersion string                     `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	MetaData   MetaData                   `json:"metadata" yaml:"metadata"`
	Spec       HorizontalPodAutoscaleSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     HPAScaleStatus             `json:"status,omitempty" yaml:"status,omitempty"`
}
type HorizontalPodAutoscaleSpec struct {
	Workload    ReplicaSet   `json:"workload" yaml:"workload,omitempty"`
	MinReplicas int32        `json:"minReplicas,omitempty"`
	MaxReplicas int32        `json:"maxReplicas"`
	Metrics     []MetricSpec `json:"metrics,omitempty"`
}
type MetricSpec struct {
	Behavior *HorizontalPodAutoscaleBehavior `json:"behavior,omitempty"`
}

type HorizontalPodAutoscaleBehavior struct {
	ScaleUp *HPAScalingRules `json:"scaleUp,omitempty"`
}
type HPAScalingRules struct {
}
type HPAScaleStatus struct {
}
