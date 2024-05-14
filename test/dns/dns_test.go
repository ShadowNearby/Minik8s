package test

import (
	"context"
	"encoding/json"
	"minik8s/config"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	"minik8s/utils"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDNS(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	client := storage.CreateEtcdStorage(config.DefaultEtcdEndpoints)
	if client == nil {
		t.Fatalf("error in create storage")
	}
	ctx := context.Background()
	content, err := os.ReadFile("dns_records.json")
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	dnsRecord := core.DNSRecord{}
	json.Unmarshal(content, &dnsRecord)
	dnsRecords := []core.DNSRecord{}
	dnsRecords = append(dnsRecords, dnsRecord)
	utils.GenerateNginxFile(dnsRecords)
	dnsKey := utils.GenerateDNSPath(dnsRecord.Host)
	entry := core.DNSEntry{Host: config.NginxListenIP}
	err = client.Put(ctx, dnsKey, entry)
	if err != nil {
		t.Errorf("error put in storage")
	}

	err = utils.StartNginx()
	if err != nil {
		utils.StopNginx()
		t.Fail()
	}
	path := dnsRecord.Paths[0]
	err = utils.CreateHelloServer(path.Port, 0)
	if err != nil {
		t.Errorf("error in create hello")
	}
	res, err := utils.TestHelloServer(dnsRecord.Host, 0)
	if err != nil || res != true {
		t.Errorf("error in test hello")
	}

	err = utils.DeleteHelloServer(path.Port, 0)
	if err != nil {
		t.Errorf("error in delete hello")
	}
	err = storage.Del(dnsKey)
	if err != nil {
		t.Errorf("error in delete storage")
	}
	err = utils.StopNginx()
	if err != nil {
		t.Fail()
	}
}
