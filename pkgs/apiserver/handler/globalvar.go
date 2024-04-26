package handler

import "minik8s/pkgs/apiserver/storage"

var etcdClient = storage.CreateEtcdStorage(storage.DefaultEndpoints)
