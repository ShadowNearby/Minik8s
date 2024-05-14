package config

const (
	LocalServerIp   = "127.0.0.1"
	clusterMasterIP = "0.0.0.0"

	ApiServerPort  = "8090"
	etcdServerPort = "2380"
	clusterMode    = true
)

const PodCIDR = "10.244.0.0/16"
const NodePort = "10250"

var DefaultEtcdEndpoints = []string{"localhost:2380"}

type apiSpace string

func GetMasterIp() string {
	if clusterMode {
		return clusterMasterIP
	} else {
		return LocalServerIp
	}

}
