package config

const (
	LocalServerIp      = "127.0.0.1"
	FunctionServerPort = "8081"
	ApiServerPort      = "8090"
	etcdServerPort     = "2380"
	clusterMode        = true
	FunctionPod        = "5000"
)

var ClusterMasterIP = "172.31.184.139"

const PodCIDR = "10.244.0.0/16"
const NodePort = "10250"

var DefaultEtcdEndpoints = []string{"localhost:2380"}
