package config

import "fmt"

const DNSPathPrefix = "/dnspath"

var NginxListenIP = "172.31.184.139"

var NginxListenAddr = fmt.Sprintf("%s:%d", NginxListenIP, 80)

const ContainerResolvPath = "/etc/resolv.conf"

const TempResolvPath = "/tmp/resolv.conf"

const ContainerHostsPath = "/etc/hosts"

const TempHostsPath = "/tmp/hosts"
