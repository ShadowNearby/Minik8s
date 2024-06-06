package core

import (
	"encoding/json"
)

type Volume struct {
	Name                  string                `json:"name" yaml:"name"`
	Nfs                   NfsVolumeAttributes   `json:"nfs,omitempty" yaml:"nfs,omitempty"`
	PersistentVolumeClaim PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty" yaml:"persistentVolumeClaim,omitempty"`
}

type PersistentVolume struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      PersistentVolumeSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status    PersistentVolumeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

func (p PersistentVolume) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

type PersistentVolumeSpec struct {
	Capacity                      ResourcesConfig               `json:"capacity" yaml:"capacity"`
	AccessMode                    PersistentVolumeAccessMode    `json:"accessMode" yaml:"accessMode"`
	PersistentVolumeReclaimPolicy PersistentVolumeReclaimPolicy `json:"persistentVolumeReclaimPolicy" yaml:"persistentVolumeReclaimPolicy"`
	StorageClassName              string                        `json:"storageClassName" yaml:"storageClassName"`
	Nfs                           NfsVolumeAttributes           `json:"nfs" yaml:"nfs"`
}

type NfsVolumeAttributes struct {
	Server           string `json:"server" yaml:"server"`
	Share            string `json:"share" yaml:"share"`
	MountPermissions string `json:"mountPermissions,omitempty" yaml:"mountPermissions,omitempty"`
	ReadOnly         bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type PersistentVolumeStatus struct {
	Phase PersistentVolumePhase `json:"phase" yaml:"phase"`
}

type PersistentVolumeClaim struct {
	BasicInfo `json:",inline" yaml:",inline"`
	Spec      PersistentVolumeClaimSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status    PersistentVolumeClaimStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

func (p PersistentVolumeClaim) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
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

type PersistentVolumeReclaimPolicy string

const (
	Retain  PersistentVolumeReclaimPolicy = "Retain"
	Recycle PersistentVolumeReclaimPolicy = "Recycle"
	Delete  PersistentVolumeReclaimPolicy = "Delete"
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
