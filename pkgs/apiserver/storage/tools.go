package storage

import (
	"encoding/json"
	"minik8s/config"
	"reflect"
	"sync"

	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
)

var etcdClient = CreateEtcdStorage(config.DefaultEtcdEndpoints)
var RedisInstance = &Redis{
	Client:   CreateRedisClient(),
	Channels: make(map[string]*redis.PubSub),
}
var FunctionInstance = &Redis{
	Client:   createFunctionClient(),
	Channels: make(map[string]*redis.PubSub),
}
var storageLock sync.Mutex

func Put(key string, val any) error {
	storageLock.Lock()
	err := RedisInstance.redisSet(key, val)
	if err != nil {
		logger.Errorf("redis cannot put: %s, error: %s", key, err.Error())
		//return err
	}
	err = TaskQueue.Enqueue(func() {
		err := etcdClient.Put(ctx, key, val)
		if err != nil {
			logger.Errorf("etcd cannot put: %s", key)
		}
	})
	if err != nil {
		logger.Errorf("cannot assign task")
		storageLock.Unlock()
		return err
	}
	storageLock.Unlock()
	return nil
}

func Get(key string, ptr any) error {
	storageLock.Lock()
	err := RedisInstance.redisGet(key, ptr)
	if err == nil {
		storageLock.Unlock()
		return err
	}
	err = etcdClient.Get(ctx, key, ptr)
	if err != nil {
		storageLock.Unlock()
		return err
	}
	storageLock.Unlock()
	return nil
}

func Del(keys ...string) error {
	storageLock.Lock()
	err := RedisInstance.redisDel(keys...)
	if err != nil {
		storageLock.Unlock()
		return err
	}
	err = TaskQueue.Enqueue(func() {
		for _, key := range keys {
			err := etcdClient.Delete(ctx, key)
			if err != nil {
				logger.Errorf("etcd del failed: %s", err.Error())
			}
		}
	})
	if err != nil {
		storageLock.Unlock()
		return err
	}
	storageLock.Unlock()
	return nil
}

func RangeGet(prefix string, ptr any) error {
	storageLock.Lock()
	var err error
	res, err := RedisInstance.redisRangeOp(prefix, OpGet)
	if err != nil {
		res, err = etcdClient.EtcdRangeOp(prefix, OpGet)
		if err != nil {
			storageLock.Unlock()
			return err
		}
	}
	listType := reflect.TypeOf(ptr).Elem()
	newVal := reflect.MakeSlice(listType, len(res), len(res))
	for i, item := range res {
		resValue := reflect.New(listType.Elem())
		err := json.Unmarshal([]byte(item.(string)), resValue.Interface())
		if err != nil {
			logger.Errorf("cannot unmarshal: %s", err.Error())
			storageLock.Unlock()
			return err
		}
		newVal.Index(i).Set(resValue.Elem())
	}
	reflect.ValueOf(ptr).Elem().Set(newVal)
	storageLock.Unlock()
	return nil
}

func RangeDel(prefix string) error {
	storageLock.Lock()
	// write redis first
	_, err := RedisInstance.redisRangeOp(prefix, OpDel)
	if err != nil {
		logger.Errorf("cannot del in redis")
	}
	// del in etcd in bg
	err = TaskQueue.Enqueue(func() {
		_, err := etcdClient.EtcdRangeOp(prefix, OpDel)
		if err != nil {
			logger.Errorf("cannot del in etcd")
		}
	})
	if err != nil {
		storageLock.Unlock()
		return err
	}
	storageLock.Unlock()
	return nil
}
