package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
	"reflect"
)

var etcdClient = CreateEtcdStorage(DefaultEndpoints)
var RedisInstance = &Redis{
	Client:   createRedisClient(),
	Channels: make(map[string]*redis.PubSub),
}

func Put(key string, val any) error {
	err := RedisInstance.RedisSet(key, val)
	if err != nil {
		logger.Errorf("redis cannot put: %s", key)
		//return err
	}
	err = TaskQueue.Enqueue(func() {
		err := etcdClient.Put(context.Background(), key, val)
		if err != nil {
			logger.Errorf("etcd cannot put: %s", key)
		}
	})
	if err != nil {
		logger.Errorf("cannot assign task")
		return err
	}
	return nil
}

func Get(key string) (val any) {
	err := RedisInstance.RedisGet(key, &val)
	if err == nil {
		return
	}
	err = etcdClient.Get(ctx, key, &val)
	if err != nil {
		return nil
	}
	return
}

func Del(keys ...string) error {
	err := RedisInstance.redisDel(keys...)
	if err != nil {
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
		return err
	}
	return nil
}

func RangeGet(prefix string, val any) error {
	// redis first
	var err error
	res, err := RedisInstance.redisRangeOp(prefix, OpGet)
	if err != nil {
		// etcd next
		res, err = etcdClient.EtcdRangeOp(prefix, OpGet)
		if err != nil {
			return err
		}
	}
	valType := reflect.TypeOf(val).Elem().Elem()
	newVal := reflect.MakeSlice(valType, len(res), len(res))
	for i, item := range res {
		resValue := reflect.ValueOf(item)
		newVal.Index(i).Set(resValue)
	}

	reflect.ValueOf(val).Elem().Set(newVal)
	return nil
}

func RangeDel(prefix string) error {
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
		return err
	}
	return nil
}
