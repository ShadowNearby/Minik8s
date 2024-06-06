package handler

import (
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func DNSKeyPrefix(namespace string, name string) string {

	return fmt.Sprintf("/dns/object/%s/%s", namespace, name)
}

func DNSListKeyPrefix(namespace string) string {
	return fmt.Sprintf("/dns/object/%s", namespace)
}

func DNSAllKeyPrefix() string {
	return "/dns/object"
}

// CreateDNSHandler POST /api/v1/namespaces/:namespace/dns
func CreateDNSHandler(c *gin.Context) {
	var dnsRecord core.DNSRecord
	err := c.Bind(&dnsRecord)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs dns config type"})
		return
	}
	if dnsRecord.MetaData.Namespace == "" {
		dnsRecord.MetaData.Namespace = "default"
	}
	for i, path := range dnsRecord.Paths {
		if path.IP == "" || path.Port == 0 {
			ip, port, err := utils.GetServiceAddress(dnsRecord.MetaData.Namespace, path.Service)
			if err != nil {
				logrus.Errorf("get service %s:%s addr error", dnsRecord.MetaData.Namespace, path.Service)
			}
			dnsRecord.Paths[i].IP = ip
			dnsRecord.Paths[i].Port = port
		}
	}
	err = storage.Put(DNSKeyPrefix(dnsRecord.MetaData.Namespace, dnsRecord.MetaData.Name), dnsRecord)
	if err != nil {
		logrus.Error("error in create dns in storage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	path := utils.GenerateDNSPath(dnsRecord.Host)
	if path == "" {
		logrus.Errorf("generate dns path error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate dns path"})
		return
	}

	err = storage.Put(path, core.DNSEntry{Host: config.NginxListenIP})
	if err != nil {
		logrus.Errorf("error put DNS path: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot put dns path"})
		return
	}

	err = utils.UpdateNginx()
	if err != nil {
		logrus.Errorf("error update nginx %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update nginx"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetDNSHandler GET /api/v1/namespaces/:namespace/dns/:name
func GetDNSHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if len(name) == 0 || len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace and name"})
		return
	}
	dnsRecord := core.DNSRecord{}
	err := storage.Get(DNSKeyPrefix(namespace, name), &dnsRecord)
	if err != nil {
		logrus.Errorf("error get dns record err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get all records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(dnsRecord)})
}

// GetDNSListHandler GET /api/v1/namespaces/:namespace/dns
func GetDNSListHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	if len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace"})
		return
	}
	dnsRecords := []core.DNSRecord{}
	err := storage.RangeGet(DNSListKeyPrefix(namespace), &dnsRecords)
	if err != nil {
		logrus.Errorf("error get dns records err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get all records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(dnsRecords)})
}

// GetAllDNSHandler GET /api/v1/dns
func GetAllDNSHandler(c *gin.Context) {
	dnsRecords := []core.DNSRecord{}
	err := storage.RangeGet(DNSAllKeyPrefix(), &dnsRecords)
	if err != nil {
		logrus.Errorf("error get all dns records err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get all records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": utils.JsonMarshal(dnsRecords)})
}

// DeleteDNSHandler DELETE /api/v1/namespaces/:namespace/dns/:name
func DeleteDNSHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if len(name) == 0 || len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace and name"})
		return
	}

	dnsRecord := core.DNSRecord{}
	err := storage.Get(DNSKeyPrefix(namespace, name), &dnsRecord)
	if err != nil {
		logrus.Error("error in get dns in storage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}

	path := utils.GenerateDNSPath(dnsRecord.Host)
	if path == "" {
		logrus.Errorf("generate dns path error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate dns path"})
		return
	}
	err = storage.Del(path)
	if err != nil {
		logrus.Errorf("error del dns path: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del dns path"})
		return
	}

	err = storage.Del(DNSKeyPrefix(namespace, name))
	if err != nil {
		logrus.Errorf("error del dns object: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del dns object"})
		return
	}

	err = utils.UpdateNginx()
	if err != nil {
		logrus.Errorf("error update nginx %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update nginx"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// UpdateDNSHandler PUT /api/v1/namespaces/:namespace/dns/:name
func UpdateDNSHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if len(name) == 0 || len(namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs namespace and name"})
		return
	}

	oldRecord := core.DNSRecord{}
	err := storage.Get(DNSKeyPrefix(namespace, name), &oldRecord)
	if err != nil {
		logrus.Error("error in get dns in storage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get data"})
		return
	}

	oldPath := utils.GenerateDNSPath(oldRecord.Host)
	if oldPath == "" {
		logrus.Errorf("generate dns path error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate dns path"})
		return
	}
	err = storage.Del(oldPath)
	if err != nil {
		logrus.Errorf("error del dns path: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del dns path"})
		return
	}

	err = storage.Del(DNSKeyPrefix(namespace, name))
	if err != nil {
		logrus.Errorf("error del dns object: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot del dns object"})
		return
	}

	var dnsRecord core.DNSRecord
	err = c.Bind(&dnsRecord)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs dns config type"})
		return
	}

	if dnsRecord.MetaData.Namespace == "" {
		dnsRecord.MetaData.Namespace = "default"
	}
	for i, path := range dnsRecord.Paths {
		if path.IP == "" || path.Port == 0 {
			ip, port, err := utils.GetServiceAddress(dnsRecord.MetaData.Namespace, path.Service)
			if err != nil {
				logrus.Errorf("get service %s:%s addr error", dnsRecord.MetaData.Namespace, path.Service)
			}
			dnsRecord.Paths[i].IP = ip
			dnsRecord.Paths[i].Port = port
		}
	}
	err = storage.Put(DNSKeyPrefix(dnsRecord.MetaData.Namespace, dnsRecord.MetaData.Name), dnsRecord)
	if err != nil {
		logrus.Error("error in create dns in storage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
	path := utils.GenerateDNSPath(dnsRecord.Host)
	if path == "" {
		logrus.Errorf("generate dns path error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate dns path"})
		return
	}

	err = storage.Put(path, core.DNSEntry{Host: config.NginxListenIP})
	if err != nil {
		logrus.Errorf("error put DNS path: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot put dns path"})
		return
	}

	err = utils.UpdateNginx()
	if err != nil {
		logrus.Errorf("error update nginx %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update nginx"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
