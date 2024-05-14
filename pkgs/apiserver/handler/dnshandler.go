package handler

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
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

// CreateDNSHandler POST /api/v1/namespaces/:namespace/dns
func CreateDNSHandler(c *gin.Context) {
	var dns core.DNSRecord
	err := c.Bind(&dns)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needs dns config type"})
		return
	}
	if dns.MetaData.Namespace == "" {
		dns.MetaData.Namespace = "default"
	}
	key := DNSKeyPrefix(dns.MetaData.Namespace, dns.MetaData.Name)
	err = storage.Put(key, dns)
	if err != nil {
		logrus.Error("error in create dns in storage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot store data"})
		return
	}
}

// GetDNSHandler GET /api/v1/namespaces/:namespace/dns/:name
func GetDNSHandler(c *gin.Context) {}

// GetDNSListHandler GET /api/v1/namespaces/:namespace/dns
func GetDNSListHandler(c *gin.Context) {}

// DeleteDNSHandler DELETE /api/v1/namespaces/:namespace/dns/:name
func DeleteDNSHandler(c *gin.Context) {}

// UpdateDNSHandler PUT /api/v1/namespaces/:namespace/dns/:name
func UpdateDNSHandler(c *gin.Context) {}
