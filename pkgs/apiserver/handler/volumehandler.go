package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/pkgs/volume"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VolumeKeyPrefix(namespace string, name string) string {
	return fmt.Sprintf("/volumes/object/%s/%s", namespace, name)
}

func VolumeListKeyPrefix(namespace string) string {
	return fmt.Sprintf("/volumes/object/%s", namespace)
}

func CsiVolumeKeyPrefix(name string) string {
	return fmt.Sprintf("/csivolumes/object/%s", name)
}

func CsiVolumeListKeyPrefix() string {
	return "/csivolumes/object"
}

// CreateVolume POST /api/v1/namespaces/:namespace/volumes
func CreateVolume(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need namespace"})
		return
	}
	pv := &core.PersistentVolume{}
	if err := c.Bind(pv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := VolumeKeyPrefix(namespace, pv.MetaData.Name)
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
