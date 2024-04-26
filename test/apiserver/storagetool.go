package apiserver

import (
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/apiserver/storage"
)

func ToolTest() {
	go bgTask()
	err := storage.Put("test1", "haha")
	if err != nil {
		logger.Errorf("put error: %s", err.Error())
		return
	}
	logger.Infof("pass get")
	var val string
	err = storage.Get("test1", &val)
	if err != nil {
		logger.Errorf("get error: %s", err.Error())
	}
	logger.Printf("get test value: %s", val)
	logger.Printf("pass get")
	_ = storage.Put("test2", "hahaha")
	var strs = make([]string, 2)
	logger.Printf("start range get")
	err = storage.RangeGet("test", &strs)
	if err != nil {
		logger.Errorf("range get error: %s", err.Error())
		return
	}
	logger.Printf("len: %d", len(strs))
	err = storage.RangeDel("test")
	if err != nil {
		logger.Errorf("range del error: %s", err.Error())
		return
	}
}

func bgTask() {
	for {
		if storage.TaskQueue.GetLen() > 0 {
			value, err := storage.TaskQueue.Dequeue()
			if err != nil {
				logger.Errorf("get bg task error: %s", err.Error())
				continue
			}
			task := value.(func())
			task()
		}
	}
}
