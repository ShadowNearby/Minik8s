package server

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
)

type APIServer struct {
	HttpServer   *gin.Engine
	EtcdStorage  *storage.EtcdStorage
	RedisStorage *storage.Redis
}

func (s *APIServer) Run(addr string) error {
	err := s.HttpServer.Run(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	s.RedisStorage = storage.RedisInstance
	s.RedisStorage.InitChannels()
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
	return &APIServer{HttpServer: gin.Default(), EtcdStorage: s}
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
