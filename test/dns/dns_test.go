package test

import (
	"context"
	"encoding/json"
	"fmt"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestDNSBasic(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	client := storage.CreateEtcdStorage(config.DefaultEtcdEndpoints)
	if client == nil {
		t.Fatalf("error in create storage")
	}
	ctx := context.Background()
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "dns/dns_records.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	dnsRecord := core.DNSRecord{}
	err = json.Unmarshal(content, &dnsRecord)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	dnsRecords := []core.DNSRecord{}
	dnsRecords = append(dnsRecords, dnsRecord)
	err = utils.GenerateNginxFile(dnsRecords)
	if err != nil {
		t.Errorf("error in generate nginx file")
	}
	dnsKey := utils.GenerateDNSPath(dnsRecord.Host)
	entry := core.DNSEntry{Host: config.NginxListenIP}
	err = client.Put(ctx, dnsKey, entry)
	if err != nil {
		t.Errorf("error put in storage")
	}

	err = utils.ReloadNginx()
	if err != nil {
		t.Fail()
	}
	path := dnsRecord.Paths[0]
	err = utils.CreateHelloServer(path.Port, 0)
	if err != nil {
		t.Errorf("error in create hello")
	}
	time.Sleep(3 * time.Second)
	res, err := utils.TestHelloServer(fmt.Sprintf("%s/%s", dnsRecord.Host, path.Path), 0)
	if err != nil || res != true {
		t.Errorf("error in test hello")
	}

	err = utils.DeleteHelloServer(path.Port, 0)
	if err != nil {
		t.Errorf("error in delete hello")
	}
	err = client.Delete(ctx, dnsKey)
	if err != nil {
		t.Errorf("error in delete storage")
	}
}

func TestDNSApi(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", utils.ExamplePath, "dns/dns_records.json"))
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	dnsRecord := core.DNSRecord{}
	err = json.Unmarshal(content, &dnsRecord)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	err = utils.CreateObject(core.ObjDNS, dnsRecord.MetaData.Namespace, dnsRecord)
	if err != nil {
		t.Errorf("error in create dns err: %s", err.Error())
	}

	path := dnsRecord.Paths[0]
	err = utils.CreateHelloServer(path.Port, 0)
	if err != nil {
		t.Errorf("error in create hello")
	}
	time.Sleep(3 * time.Second)
	res, err := utils.TestHelloServer(fmt.Sprintf("%s/%s", dnsRecord.Host, path.Path), 0)
	if err != nil || res != true {
		t.Errorf("error in test hello")
	}

	err = utils.DeleteHelloServer(path.Port, 0)
	if err != nil {
		t.Errorf("error in delete hello")
	}

	err = utils.DeleteObject(core.ObjDNS, dnsRecord.MetaData.Namespace, dnsRecord.MetaData.Name)
	if err != nil {
		t.Errorf("error in del dns err: %s", err.Error())
	}
}
