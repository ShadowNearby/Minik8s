package storage

import (
	logger "github.com/sirupsen/logrus"
	"testing"
)

func TestTools(t *testing.T) {
	go bgTask()
	err := Put("test1", "haha")
	if err != nil {
		t.Errorf("put error: %s", err.Error())
		return
	}
	var val string
	err = Get("test1", &val)
	if err != nil {
		t.Errorf("get error: %s", err.Error())
	}
	logger.Infof("get test value: %s", val)
	_ = Put("test2", "hahaha")
	var strs []string
	err = RangeGet("test", &strs)
	if err != nil {
		t.Errorf("range get error: %s", err.Error())
		return
	}
	err = RangeDel("test")
	if err != nil {
		t.Errorf("range del error: %s", err.Error())
		return
	}
}

func bgTask() {
	for {
		if TaskQueue.GetLen() > 0 {
			value, err := TaskQueue.Dequeue()
			if err != nil {
				logger.Errorf("get bg task error: %s", err.Error())
				continue
			}
			task := value.(func())
			task()
		}
	}
}
