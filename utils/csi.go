package utils

import (
	core "minik8s/pkgs/apiobject"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func AsCSIAccessMode(am core.PersistentVolumeAccessMode) csi.VolumeCapability_AccessMode_Mode {
	switch am {
	case core.ReadWriteOnce:
		return csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER
	case core.ReadOnlyMany:
		return csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY
	case core.ReadWriteMany:
		return csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER
	// This mapping exists to enable CSI drivers that lack the
	// SINGLE_NODE_MULTI_WRITER capability to work with the
	// ReadWriteOncePod access mode.
	case core.ReadWriteOncePod:
		return csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER
	}
	return csi.VolumeCapability_AccessMode_UNKNOWN
}

func AsCSIReadOnly(am core.PersistentVolumeAccessMode) bool {
	return am == core.ReadOnlyMany
}
