package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	LocalServerIp        = "127.0.0.1"
	ApiServerPort        = "8090"
	RegistryPort         = "5000"
	ClusterMasterIP      = "192.168.1.12"
	PodCIDR              = "10.244.0.0/16"
	NodePort             = "10250"
	DefaultEtcdEndpoints = []string{"localhost:2380"}
	HeartbeatInterval    = 15 * time.Second
)

var (
	CsiSockAddr          = "/run/csi/csi.sock"
	CsiStagingTargetPath = "/mnt/staging"
	CsiMntPath           = "/mnt/minik8s"
	CsiServerIP          = ClusterMasterIP
	CsiStorageClassName  = "nfs-csi"
)

var (
	DNSPathPrefix       = "/dnspath"
	NginxListenIP       = ClusterMasterIP
	NginxListenAddr     = fmt.Sprintf("%s:%d", NginxListenIP, 80)
	NginxStarted        = false
	ContainerResolvPath = "/etc/resolv.conf"
	TempResolvPath      = "/tmp/resolv.conf"
	ContainerHostsPath  = "/etc/hosts"
	TempHostsPath       = "/tmp/hosts"
)

var (
	PrometheusNodeFilePath   = "prometheus/sd_node.json"
	PrometheusPodFilePath    = "prometheus/sd_pod.json"
	PrometheusScrapeInterval = 10 * time.Second
)

var (
	FunctionRetryTimes        = 10
	FunctionServerIp          = "master"
	FunctionThreshold         = 6
	FunctionConnectTime       = 30 * time.Second
	ServerlessScaleToZeroTime = 90 * time.Second
	ServerlessPort            = "18080"
)

var (
	HPAScaleInterval = 30 * time.Second
)

func InitConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		logrus.Errorf("parse config error %s", err.Error())
		return err
	}

	LocalServerIp = cfg.LocalServerIp
	ApiServerPort = cfg.ApiServerPort
	ClusterMasterIP = cfg.ClusterMasterIP
	PodCIDR = cfg.PodCIDR
	NodePort = cfg.NodePort
	DefaultEtcdEndpoints = cfg.DefaultEtcdEndpoints
	HeartbeatInterval, err = time.ParseDuration(cfg.HeartbeatInterval)
	if err != nil {
		logrus.Errorf("error in parse HeartbeatInterval %s", err.Error())
		return err
	}
	CsiSockAddr = cfg.CsiSockAddr
	CsiStagingTargetPath = cfg.CsiStagingTargetPath
	CsiMntPath = cfg.CsiMntPath
	CsiServerIP = cfg.CsiServerIP
	CsiStorageClassName = cfg.CsiStorageClassName
	DNSPathPrefix = cfg.DNSPathPrefix
	NginxListenIP = cfg.NginxListenIP
	NginxListenAddr = cfg.NginxListenIP
	NginxStarted = cfg.NginxStarted
	ContainerResolvPath = cfg.ContainerResolvPath
	TempResolvPath = cfg.TempResolvPath
	ContainerHostsPath = cfg.ContainerHostsPath
	TempHostsPath = cfg.TempHostsPath
	PrometheusNodeFilePath = cfg.PrometheusNodeFilePath
	PrometheusPodFilePath = cfg.PrometheusPodFilePath

	PrometheusScrapeInterval, err = time.ParseDuration(cfg.PrometheusScrapeInterval)
	if err != nil {
		logrus.Errorf("error in parse PrometheusScrapeInterval %s", err.Error())
		return err
	}
	FunctionRetryTimes = cfg.FunctionRetryTimes
	FunctionServerIp = cfg.FunctionServerIp
	FunctionThreshold = cfg.FunctionThreshold
	FunctionConnectTime, err = time.ParseDuration(cfg.FunctionConnectTime)
	if err != nil {
		logrus.Errorf("error in parse FunctionConnectTime %s", err.Error())
		return err
	}
	ServerlessScaleToZeroTime, err = time.ParseDuration(cfg.ServerlessScaleToZeroTime)
	if err != nil {
		logrus.Errorf("error in parse ServerlessScaleToZeroTime %s", err.Error())
		return err

	}
	ServerlessPort = cfg.ServelessIP
	HPAScaleInterval, err = time.ParseDuration(cfg.HPAScaleInterval)
	if err != nil {
		logrus.Errorf("error in parse HPAScaleInterval %s", err.Error())
		return err
	}

	logrus.Info("Configuration parsed successfully")
	return nil
}

type Config struct {
	LocalServerIp             string   `json:"LocalServerIp"`
	ApiServerPort             string   `json:"ApiServerPort"`
	ClusterMasterIP           string   `json:"ClusterMasterIP"`
	PodCIDR                   string   `json:"PodCIDR"`
	NodePort                  string   `json:"NodePort"`
	DefaultEtcdEndpoints      []string `json:"DefaultEtcdEndpoints"`
	HeartbeatInterval         string   `json:"HeartbeatInterval"`
	CsiSockAddr               string   `json:"CsiSockAddr"`
	CsiStagingTargetPath      string   `json:"CsiStagingTargetPath"`
	CsiMntPath                string   `json:"CsiMntPath"`
	CsiServerIP               string   `json:"CsiServerIP"`
	CsiStorageClassName       string   `json:"CsiStorageClassName"`
	DNSPathPrefix             string   `json:"DNSPathPrefix"`
	NginxListenIP             string   `json:"NginxListenIP"`
	NginxStarted              bool     `json:"NginxStarted"`
	ContainerResolvPath       string   `json:"ContainerResolvPath"`
	TempResolvPath            string   `json:"TempResolvPath"`
	ContainerHostsPath        string   `json:"ContainerHostsPath"`
	TempHostsPath             string   `json:"TempHostsPath"`
	PrometheusNodeFilePath    string   `json:"PrometheusNodeFilePath"`
	PrometheusPodFilePath     string   `json:"PrometheusPodFilePath"`
	PrometheusScrapeInterval  string   `json:"PrometheusScrapeInterval"`
	FunctionRetryTimes        int      `json:"FunctionRetryTimes"`
	FunctionServerIp          string   `json:"FunctionServerIp"`
	FunctionThreshold         int      `json:"FunctionThreshold"`
	FunctionConnectTime       string   `json:"FunctionConnectTime"`
	ServerlessScaleToZeroTime string   `json:"ServerlessScaleToZeroTime"`
	ServelessIP               string   `json:"ServerlessPort"`
	HPAScaleInterval          string   `json:"HPAScaleInterval"`
}
