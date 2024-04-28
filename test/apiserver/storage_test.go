package apiserver

import (
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/apiserver/storage"
	"sync"
	"testing"
)

func TestToolTest(t *testing.T) {
	var wg sync.WaitGroup
	go bgTask()
	wg.Add(2)
	go func(g *sync.WaitGroup) {
		defer wg.Done()
		err := storage.Put("test1", "haha")
		if err != nil {
			t.Errorf("put error: %s", err.Error())
			return
		}
		var val string
		err = storage.Get("test1", &val)
		if err != nil {
			t.Errorf("get error: %s", err.Error())
		}
		if val != "haha" {
			t.Errorf("get value wrong")
		}
	}(&wg)
	go func(g *sync.WaitGroup) {
		defer wg.Done()
		_ = storage.Put("test2", "hahaha")
	}(&wg)
	wg.Wait()
	var strs = make([]string, 2)
	logger.Printf("start range get")
	err := storage.RangeGet("test", &strs)
	if err != nil {
		t.Errorf("range get error: %s", err.Error())
		return
	}
	if "hahaha" == strs[0] {
		if strs[1] != "haha" {
			t.Errorf("range get error")
		}
	} else if strs[0] == "haha" {
		if strs[1] != "hahaha" {
			t.Errorf("range get error")
		}
	} else {
		t.Errorf("range get error")
	}
	err = storage.RangeDel("test")
	if err != nil {
		t.Errorf("range del error: %s", err.Error())
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
