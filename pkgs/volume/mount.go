package volume

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/utils"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
)

func HandleDynamicVolumes(mounts []core.Volume) error {
	for _, volume := range mounts {
		pv := &core.PersistentVolume{}
		pv.MetaData.Name = volume.Name
		pv.Spec.Nfs = volume.Nfs
		err := utils.CreateObjectWONamespace(core.ObjVolume, pv)
		if err != nil {
			logrus.Errorf("error in create volume: %s", err.Error())
			return err
		}
	}
	return nil
}

func HandleMount(mounts []core.VolumeMountConfig) error {
	for _, mount := range mounts {
		if mount.Name != "" {
			pv := &core.PersistentVolume{}
			resp := utils.GetObjectWONamespace(core.ObjVolume, mount.Name)
			err := utils.JsonUnMarshal(resp, pv)
			if err != nil {
				logrus.Errorf("error in unmarshal: %s", err.Error())
				return err
			}
			csiVolume, err := CreateVolume(pv)
			if err != nil {
				logrus.Errorf("error in create volume: %s", err.Error())
				return err
			}
			mountPath := fmt.Sprintf("%s/%s", config.CsiMntPath, pv.MetaData.Name)
			err = NodePublishVolume(csiVolume.VolumeId, mountPath, pv)
			if err != nil {
				logrus.Errorf("error in publish volume: %s", err.Error())
				return err
			}
		}
	}
	return nil
}

func HandleUnmount(mounts []core.VolumeMountConfig) error {
	for _, mount := range mounts {
		if mount.Name != "" {
			volume := &csi.Volume{}
			resp := utils.GetObjectWONamespace(core.ObjCsiVolume, mount.Name)
			err := utils.JsonUnMarshal(resp, volume)
			if err != nil {
				logrus.Errorf("error in unmarshal: %s", err.Error())
				return err
			}
			pv := &core.PersistentVolume{}
			resp = utils.GetObjectWONamespace(core.ObjVolume, mount.Name)
			err = utils.JsonUnMarshal(resp, pv)
			if err != nil {
				logrus.Errorf("error in unmarshal: %s", err.Error())
				return err
			}
			mountPath := fmt.Sprintf("%s/%s", config.CsiMntPath, pv.MetaData.Name)
			err = NodeUnpublishVolume(volume.VolumeId, mountPath)
			if err != nil {
				logrus.Errorf("error in unpublish volume: %s", err.Error())
				return err
			}
			err = DeleteVolume(volume.VolumeId)
			if err != nil {
				logrus.Errorf("error in delete volume: %s", err.Error())
				return err
			}

		}
	}
	return nil
}
