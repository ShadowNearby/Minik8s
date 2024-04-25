package server

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
)

type APIServer struct {
	HttpServer  *gin.Engine
	EtcdStorage *storage.EtcdStorage
}

func (s *APIServer) Run(addr string) error {
	err := s.HttpServer.Run(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func InitNodes(storage *storage.EtcdStorage) {
	// delete all nodes' info in etcd
	key := "/registry/nodes/"
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
