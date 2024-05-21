package handler

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/volume"
	"minik8s/utils"
	"net/http"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VolumeKeyPrefix(name string) string {
	return fmt.Sprintf("/volumes/object/%s", name)
}

func CsiVolumeKeyPrefix(name string) string {
	return fmt.Sprintf("/csivolumes/object/%s", name)
}

// CreateVolumeHandler POST /api/v1/volumes
func CreateVolumeHandler(c *gin.Context) {
	pv := &core.PersistentVolume{}
	if err := c.Bind(pv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := VolumeKeyPrefix(pv.MetaData.Name)
	if err := storage.Put(key, pv); err != nil {
		logrus.Errorf("save pv error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	csiVolume, err := volume.CreateVolume(pv)
	if err != nil {
		logrus.Errorf("create csi volume error %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in create csi volume": err.Error()})
		return
	}
	key = CsiVolumeKeyPrefix(pv.MetaData.Name)
	if err := storage.Put(key, *csiVolume); err != nil {
		logrus.Errorf("save csi volume error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetVolumeHandler GET /api/v1/volumes/:name
func GetVolumeHandler(c *gin.Context) {
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs name"})
		return
	}
	pv := &core.PersistentVolume{}
	key := VolumeKeyPrefix(name)
	err := storage.Get(key, &pv)
	if err != nil {
		logrus.Errorf("error in get pv %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(*pv)})
}

// GetCsiVolumeHandler GET /api/v1/csivolumes/:name
func GetCsiVolumeHandler(c *gin.Context) {
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs name"})
		return
	}
	volume := &csi.Volume{}
	key := VolumeKeyPrefix(name)
	err := storage.Get(key, &volume)
	if err != nil {
		logrus.Errorf("error in get csi volume %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(*volume)})
}

// GetVolumeHandler DELETE /api/v1/volumes/:name
func DeleteVolumeHandler(c *gin.Context) {
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs name"})
		return
	}
	csiVolume := &csi.Volume{}
	objVolume := &core.PersistentVolume{}
	key := CsiVolumeKeyPrefix(name)
	err := storage.Get(key, csiVolume)

	if err != nil {
		logrus.Errorf("error in get csi volume %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = storage.Del(key)
	if err != nil {
		logrus.Errorf("error in del csi volume %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key = VolumeKeyPrefix(name)

	err = storage.Get(key, objVolume)

	if err != nil {
		logrus.Errorf("error in get csi volume %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = storage.Del(key)
	if err != nil {
		logrus.Errorf("error in del pv %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mountPath := fmt.Sprintf("%s/%s", config.CsiMntPath, objVolume.MetaData.Name)
	err = volume.NodeUnpublishVolume(csiVolume.VolumeId, mountPath)

	if err != nil {
		logrus.Errorf("error in unmount volume %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = volume.DeleteVolume(csiVolume.VolumeId)

	if err != nil {
		logrus.Errorf("error in delete volume %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
