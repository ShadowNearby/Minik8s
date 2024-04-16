package storage

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

type EtcdStorage struct {
	client *clientv3.Client
}

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

func (e *EtcdStorage) Put(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	res, err := e.client.Put(ctx, key, string(jsonValue))
	if err != nil {
		return err
	}
	log.Infoln(res.Header.String())
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
