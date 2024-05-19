package core

type Volume struct {
	Name                  string                            `json:"name" yaml:"name"`
	PersistentVolumeClaim PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty" yaml:"persistentVolumeClaim,omitempty"`
}

type PersistentVolumeClaimVolumeSource struct {
	ClaimName string `json:"claimName" yaml:"claimName"`
	ReadOnly  bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type PersistentVolume struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      PersistentVolumeSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status    PersistentVolumeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type PersistentVolumeSpec struct {
	Capacity                      ResourcesConfig            `json:"capacity" yaml:"capacity"`
	AccessMode                    PersistentVolumeAccessMode `json:"accessMode" yaml:"accessMode"`
	PersistentVolumeReclaimPolicy string                     `json:"persistentVolumeReclaimPolicy" yaml:"persistentVolumeReclaimPolicy"`
	StorageClassName              string                     `json:"storageClassName" yaml:"storageClassName"`
}

type PersistentVolumeStatus struct {
	Phase PersistentVolumePhase `json:"phase" yaml:"phase"`
}

type PersistentVolumeClaim struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      PersistentVolumeClaimSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status    PersistentVolumeClaimStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type PersistentVolumeClaimSpec struct {
	AccessMode       PersistentVolumeAccessMode `json:"accessMode" yaml:"accessMode"`
	Selector         Selector                   `json:"selector" yaml:"selector"`
	StorageClassName string                     `json:"storageClassName" yaml:"storageClassName"`
	Resources        ResourcesConfig            `json:"resources,omitempty" yaml:"resources,omitempty"`
	VolumeName       string                     `json:"volumeName,omitempty" protobuf:"bytes,3,opt,name=volumeName"`
}

type PersistentVolumeClaimStatus struct {
	Phase                     string                     `json:"phase" yaml:"phase"`
	AccessMode                PersistentVolumeAccessMode `json:"accessMode" yaml:"accessMode"`
	Capacity                  ResourceLimit              `json:"capacity" yaml:"capacity"`
	AllocatedResources        ResourceLimit              `json:"allocatedResources" yaml:"allocatedResources"`
	AllocatedResourceStatuses string                     `json:"allocatedResourceStatuses" yaml:"allocatedResourceStatuses"`
}

type PersistentVolumeAccessMode string

const (
	ReadWriteOnce    PersistentVolumeAccessMode = "ReadWriteOnce"
	ReadOnlyMany     PersistentVolumeAccessMode = "ReadOnlyMany"
	ReadWriteMany    PersistentVolumeAccessMode = "ReadWriteMany"
	ReadWriteOncePod PersistentVolumeAccessMode = "ReadWriteOncePod"
)

type PersistentVolumePhase string

const (
	VolumePending   PersistentVolumePhase = "Pending"
	VolumeAvailable PersistentVolumePhase = "Available"
	VolumeBound     PersistentVolumePhase = "Bound"
	VolumeReleased  PersistentVolumePhase = "Released"
	VolumeFailed    PersistentVolumePhase = "Failed"
)

type PersistentVolumeClaimPhase string

const (
	ClaimPending PersistentVolumeClaimPhase = "Pending"
	ClaimBound   PersistentVolumeClaimPhase = "Bound"
	ClaimLost    PersistentVolumeClaimPhase = "Lost"
)
