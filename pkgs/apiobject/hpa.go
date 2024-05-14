package core

import "time"

type HorizontalPodAutoscaler struct {
	ApiVersion string    `json:"apiVersion" yaml:"apiVersion"`
	Kind       string    `json:"kind" yaml:"kind"`
	MetaData   MetaData  `json:"metaData" yaml:"metaData"`
	Spec       HPASpec   `json:"spec" yaml:"spec"`
	Status     HPAStatus `json:"status" yaml:"status"`
}
type HPASpec struct {
	ScaleTargetRef ScaleTargetRef `json:"scaleTargetRef" yaml:"scaleTargetRef"`
	MinReplicas    int            `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas    int            `json:"maxReplicas" yaml:"maxReplicas"`
	Metrics        Metrics        `json:"metrics" yaml:"metrics"`
}

type ScaleTargetRef struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"` // support replicaset
	Name       string `json:"name" yaml:"name"`
	Namespace  string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type Metrics struct {
	Type      string     `json:"type" yaml:"type"` // "Resource"
	Resources []Resource `json:"resources" yaml:"resources"`
}

type Resource struct {
	Name   string         `json:"name" yaml:"name"` // "cpu", "memory"
	Target ResourceTarget `json:"target" yaml:"target"`
}

type ResourceTarget struct {
	Type               string `json:"type" yaml:"type"`                             // "Utilization"
	Value              uint64 `yaml:"value" json:"value"`                           // absolute value
	AverageUtilization int    `json:"averageUtilization" yaml:"averageUtilization"` // percentage
}

type HPAStatus struct {
	ObservedGeneration int       `json:"observedGeneration" yaml:"observedGeneration"`
	LastScaleTime      time.Time `json:"lastScaleTime" yaml:"lastScaleTime"`
	CurrentReplicas    int       `json:"currentReplicas" yaml:"currentReplicas"`
	DesiredReplicas    int       `json:"desiredReplicas" yaml:"desiredReplicas"`
	CurrentMetrics     []Metrics `json:"currentMetrics" yaml:"currentMetrics"`
}
