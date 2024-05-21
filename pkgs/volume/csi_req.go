package volume

import (
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
)

func CreateVolume(volume *core.PersistentVolume) (*csi.Volume, error) {
	resp, err := CsiClientInstance.ControllerClient.CreateVolume(CsiClientInstance.Context, &csi.CreateVolumeRequest{
		Name: volume.MetaData.Name,
		VolumeCapabilities: []*csi.VolumeCapability{
			{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: utils.AsCSIAccessMode(volume.Spec.AccessMode),
				},
			},
		},
		CapacityRange: &csi.CapacityRange{
			RequiredBytes: int64(volume.Spec.Capacity.Request.Storage),
			LimitBytes:    int64(volume.Spec.Capacity.Limit.Storage),
		},
		Parameters: map[string]string{
			"server": volume.Spec.Nfs.Server,
			"share":  volume.Spec.Nfs.Share,
			"subDir": "",
		},
	})
	if err != nil {
		logrus.Errorf("failed to create volume: %s", err.Error())
		return nil, err
	}
	return resp.Volume, nil
}

func DeleteVolume(volumeId string) error {
	_, err := CsiClientInstance.ControllerClient.DeleteVolume(CsiClientInstance.Context, &csi.DeleteVolumeRequest{
		VolumeId: volumeId,
	})
	if err != nil {
		logrus.Errorf("failed to get volume: %s", err.Error())
		return err
	}
	return nil
}

func NodePublishVolume(volumeId string, targetPath string, objVolume *core.PersistentVolume) error {
	_, err := CsiClientInstance.NodeClient.NodePublishVolume(CsiClientInstance.Context, &csi.NodePublishVolumeRequest{
		TargetPath: targetPath,
		VolumeId:   volumeId,
		VolumeCapability: &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: utils.AsCSIAccessMode(objVolume.Spec.AccessMode),
			},
		},
		VolumeContext: map[string]string{
			"server": objVolume.Spec.Nfs.Server,
			"share":  objVolume.Spec.Nfs.Share,
			"subdir": "",
		},
		Readonly: objVolume.Spec.Nfs.ReadOnly,
	})
	if err != nil {
		logrus.Errorf("failed to node publish volume: %s", err.Error())
		return err
	}
	return nil
}

func NodeUnpublishVolume(volumeId string, targetPath string) error {
	_, err := CsiClientInstance.NodeClient.NodeUnpublishVolume(CsiClientInstance.Context, &csi.NodeUnpublishVolumeRequest{
		TargetPath: targetPath,
		VolumeId:   volumeId,
	})
	if err != nil {
		logrus.Errorf("failed to node publish volume: %s", err.Error())
		return err
	}
	return nil
}

func ControllerPublishVolume(nodeId string, csiVolume *csi.Volume, objVolume *core.PersistentVolume) (*csi.ControllerPublishVolumeResponse, error) {
	resp, err := CsiClientInstance.ControllerClient.ControllerPublishVolume(CsiClientInstance.Context, &csi.ControllerPublishVolumeRequest{
		VolumeId: csiVolume.VolumeId,
		NodeId:   nodeId,
		VolumeCapability: &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: utils.AsCSIAccessMode(objVolume.Spec.AccessMode),
			},
		},
		Readonly:      utils.AsCSIReadOnly(objVolume.Spec.AccessMode),
		VolumeContext: csiVolume.VolumeContext,
	})
	if err != nil {
		logrus.Errorf("failed to controller publish volume: %s", err.Error())
		return nil, err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return resp, nil
}

func ControllerUnpublishVolume(csiVolume *csi.Volume, objVolume *core.PersistentVolume, nodeId string) error {
	resp, err := CsiClientInstance.ControllerClient.ControllerUnpublishVolume(CsiClientInstance.Context, &csi.ControllerUnpublishVolumeRequest{
		VolumeId: csiVolume.VolumeId,
		NodeId:   nodeId,
	})
	if err != nil {
		logrus.Errorf("failed to controller publish volume: %s", err.Error())
		return err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return nil
}

func NodeStageVolume(volumeId string, objVolume *core.PersistentVolume, controllerPublishVolumeResponse *csi.ControllerPublishVolumeResponse) error {
	resp, err := CsiClientInstance.NodeClient.NodeStageVolume(CsiClientInstance.Context, &csi.NodeStageVolumeRequest{
		VolumeId: volumeId,
		VolumeCapability: &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: utils.AsCSIAccessMode(objVolume.Spec.AccessMode),
			},
		},
		// PublishContext:    controllerPublishVolumeResponse.PublishContext,
		StagingTargetPath: config.CsiStagingTargetPath,
	})
	if err != nil {
		logrus.Errorf("failed to node stage volume: %s", err.Error())
		return err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return nil
}

func NodeUntageVolume(volumeId string, objVolume *core.PersistentVolume, controllerPublishVolumeResponse *csi.ControllerPublishVolumeResponse) error {
	resp, err := CsiClientInstance.NodeClient.NodeUnstageVolume(CsiClientInstance.Context, &csi.NodeUnstageVolumeRequest{
		VolumeId:          volumeId,
		StagingTargetPath: config.CsiStagingTargetPath,
	})
	if err != nil {
		logrus.Errorf("failed to node stage volume: %s", err.Error())
		return err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return nil
}

func ControllerGetVolume(volumeId string) (*csi.Volume, error) {
	resp, err := CsiClientInstance.ControllerClient.ControllerGetVolume(CsiClientInstance.Context, &csi.ControllerGetVolumeRequest{
		VolumeId: volumeId,
	})
	if err != nil {
		logrus.Errorf("failed to get volume: %s", err.Error())
		return nil, err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return resp.Volume, nil
}

func GetCapacity() ([]*csi.ControllerServiceCapability, error) {
	resp, err := CsiClientInstance.ControllerClient.ControllerGetCapabilities(CsiClientInstance.Context, &csi.ControllerGetCapabilitiesRequest{})
	if err != nil {
		logrus.Errorf("failed to get volume: %s", err.Error())
		return nil, err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return resp.Capabilities, nil
}

func ListVolumes() ([]*csi.ListVolumesResponse_Entry, error) {
	resp, err := CsiClientInstance.ControllerClient.ListVolumes(CsiClientInstance.Context, &csi.ListVolumesRequest{})
	if err != nil {
		logrus.Errorf("failed to list volumes: %s", err.Error())
		return nil, err
	}
	logrus.Infof("info: %s", utils.JsonMarshal(resp))
	return resp.Entries, nil
}
