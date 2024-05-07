package storage

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdStorage struct {
	client *clientv3.Client
}

var DefaultEndpoints = []string{"localhost:2380"}

func CreateEtcdStorage(endpoints []string) *EtcdStorage {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Error(err)
		return nil
	}
	return &EtcdStorage{client: client}
}

func (e *EtcdStorage) CloseEtcdStorage() error {
	return e.client.Close()
}

func (e *EtcdStorage) Get(ctx context.Context, key string, result interface{}) error {
	response, err := e.client.Get(ctx, key)
	if err != nil {
		return err
	}
	if response.Kvs == nil || len(response.Kvs) == 0 {
		err = errors.New("key not found\n")
		log.Errorf("key not found: %s \n", key)
		return err
	}
	if len(response.Kvs) > 1 {
		log.Errorf("Multi value for key: %s \n", key)
		return err
	}
	err = json.Unmarshal(response.Kvs[0].Value, result)
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStorage) GetList(ctx context.Context, prefix string, result interface{}) error {
	response, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	resultType := reflect.TypeOf(result).Elem().Elem()
	kvs := response.Kvs

	if kvs == nil || len(kvs) == 0 {
		items := reflect.MakeSlice(reflect.SliceOf(resultType), 0, 0)
		reflect.ValueOf(result).Elem().Set(items)
		return nil
	}

	items := reflect.MakeSlice(reflect.SliceOf(resultType), len(kvs), len(kvs))

	for i, kv := range kvs {
		item := reflect.New(resultType).Interface()
		if err := json.Unmarshal(kv.Value, item); err != nil {
			return err
		}
		items.Index(i).Set(reflect.ValueOf(item).Elem())
	}

	reflect.ValueOf(result).Elem().Set(items)

	return nil
}

func (e *EtcdStorage) EtcdRangeOp(prefix string, op string) ([]any, error) {
	switch op {
	case OpDel:
		resp, err := e.client.Delete(ctx, prefix, clientv3.WithPrefix())
		if err != nil {
			return nil, err
		}
		log.Infof("%d keys were deleted", resp.Deleted)
	case OpGet:
		var vals []any
		err := e.GetList(ctx, prefix, &vals)
		if err != nil {
			return nil, err
		}
		return vals, nil
	}
	return nil, errors.New("unsupported op type")
}

func (e *EtcdStorage) Put(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = e.client.Put(ctx, key, string(jsonValue))
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStorage) Delete(ctx context.Context, key string) error {
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e *EtcdStorage) Watch(ctx context.Context, key string, callback func(string, []byte) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := e.client.Watch(ctx, key, clientv3.WithPrefix())

	for {
		select {
		case response := <-ch:
			for _, event := range response.Events {
				log.Infof("[Watch] key %s value %s type %s", string(event.Kv.Key), string(event.Kv.Value), event.Type)
				err := callback(string(event.Kv.Key), event.Kv.Value)
				if err != nil {
					log.Error("watch error")
					return err
				}
			}
		case <-ctx.Done():
			log.Error("ctx done")
			return nil
		}
	}
}
