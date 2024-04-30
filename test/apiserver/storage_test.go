package test

import (
	logger "github.com/sirupsen/logrus"
	"minik8s/pkgs/apiserver/storage"
	"sync"
	"testing"
)

func TestTools(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	done := false
	go bgTask(&done, wg)
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
	logger.Infof("get test value: %s", val)
	_ = storage.Put("test2", "hahaha")
	var strs []string
	err = storage.RangeGet("test", &strs)
	if err != nil {
		t.Errorf("range get error: %s", err.Error())
		return
	}
	err = storage.RangeDel("test")
	if err != nil {
		t.Errorf("range del error: %s", err.Error())
		return
	}
	done = true
	wg.Wait()
}

func TestToolsConcurrent(t *testing.T) {
	wg1 := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	wg1.Add(1)
	done := false
	go bgTask(&done, wg1)
	wg2.Add(2)
	concurrentTask(100, "test1", "ha", wg2)
	concurrentTask(100, "test1", "haha", wg2)
	wg2.Wait()
	err := storage.Del("test1")
	if err != nil {
		logger.Errorf("del error: %s", err.Error())
		return
	}
	done = true
	wg1.Wait()
}

func TestWatch(t *testing.T) {
	wg1 := &sync.WaitGroup{}
	wg1.Add(2)
	done := false
	go bgTask(&done, wg1)
	storage.RedisInstance.CreateChannel("ch1")
	go func() {
		ch := storage.RedisInstance.SubscribeChannel("ch1")
		for msg := range ch {
			logger.Infof("recv msg: %s", msg)
		}
	}()
	go func() {
		storage.RedisInstance.PublishMessage("ch1", "hello world")
		storage.RedisInstance.PublishMessage("ch1", "hello world")
		storage.RedisInstance.PublishMessage("ch1", "hello world")
		storage.RedisInstance.PublishMessage("ch1", "hello world")
		wg1.Done()
	}()
	done = true
	wg1.Wait()
}

func bgTask(done *bool, wg *sync.WaitGroup) {
	for {
		if *done == true && storage.TaskQueue.GetLen() == 0 {
			break
		}
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
	wg.Done()
}

func concurrentTask(times int, key string, val string, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < times; i++ {
			err := storage.Put(key, val)
			if err != nil {
				logger.Fatalf("put error: %s", err.Error())
			}
			var newVal string
			err = storage.Get(key, &newVal)
			logger.Infof("times %d, key: %s, value: %s", i, key, newVal)
			if err != nil {
				logger.Fatalf("get error: %s", err.Error())
			}
		}
		defer wg.Done()
	}()

}
