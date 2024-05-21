package config

import "fmt"

const DNSPathPrefix = "/dnspath"

var NginxListenIP = ClusterMasterIP

var NginxListenAddr = fmt.Sprintf("%s:%d", NginxListenIP, 80)

var NginxStarted = false

const ContainerResolvPath = "/etc/resolv.conf"

const TempResolvPath = "/tmp/resolv.conf"

const ContainerHostsPath = "/etc/hosts"

const TempHostsPath = "/tmp/hosts"
