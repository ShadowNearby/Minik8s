package server

import (
	"context"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/handler"
	"minik8s/pkgs/apiserver/storage"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type APIServer struct {
	HttpServer   *gin.Engine
	EtcdStorage  *storage.EtcdStorage
	RedisStorage *storage.Redis
}

func (s *APIServer) Run(addr string) error {
	handler.PrometheusRegister()
	for _, route := range handler.RouteTable {
		route.Register(s.HttpServer)
	}
	go s.HttpServer.Run(addr)
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }
	go bgTask()
	select {}
}

func InitNodes(storage *storage.EtcdStorage) {
	// delete all nodes' info in etcd
	key := "/nodes/object"
	var nodes []core.Node
	err := storage.GetList(context.Background(), key, &nodes)
	if err != nil {
		log.Info("[InitNodes] the node list is empty")
	} else {
		for _, node := range nodes {
			nodeKey := key + node.NodeMetaData.Name
			err := storage.Delete(context.Background(), nodeKey)
			if err != nil {
				log.Error("[InitNodes] delete node error: ", err)
			}
		}
	}
}

func CreateAPIServer(endpoints []string) *APIServer {
	s := storage.CreateEtcdStorage(endpoints)
	if s == nil {
		return nil
	}
	storage.RedisInstance.InitChannels()
	return &APIServer{
		HttpServer:   gin.Default(),
		EtcdStorage:  s,
		RedisStorage: storage.RedisInstance,
	}
}

func bgTask() {
	for {
		if storage.TaskQueue.GetLen() > 0 {
			value, err := storage.TaskQueue.Dequeue()
			if err != nil {
				log.Errorf("get bg task error: %s", err.Error())
				continue
			}
			task := value.(func())
			task()
		}
	}
}
