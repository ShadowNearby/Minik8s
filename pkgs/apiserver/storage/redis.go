package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
	"reflect"
)

type Redis struct {
	Client   *redis.Client
	Channels map[string]*redis.PubSub
}

var ctx = context.Background()

const (
	ChannelNode    string = "NODE"
	ChannelPod     string = "POD"
	ChannelService string = "SERVICE"
	ChannelReplica string = "REPLICASET"
)

const (
	OpDel string = "delete"
	OpGet string = "get"
	OpSet string = "set"
)

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (r *Redis) redisSet(key string, value any) error {
	err := r.Client.Set(ctx, key, value, 0)
	if err != nil {
		return err.Err()
	}
	return nil
}

// redisGet bind should be a pointer
func (r *Redis) redisGet(key string, ptr any) error {
	val, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		logger.Errorf("%s does not exist", key)
	}
	if err != nil {
		return err
	}
	ptrValue := reflect.ValueOf(ptr)
	eleType := reflect.TypeOf(ptr).Elem()
	item := reflect.ValueOf(val)
	if !item.Type().AssignableTo(eleType) {
		return errors.New("value type does not match pointer type")
	}
	ptrValue.Elem().Set(item)
	return nil
}

func (r *Redis) redisDel(keys ...string) error {
	_, err := r.Client.Del(ctx, keys...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) redisRangeOp(prefix string, op string) ([]any, error) {
	logger.Printf("in redis range op")
	var cursor uint64 = 0
	var keys []string
	var vals []any
	var err error
	keys, cursor, err = r.Client.Scan(ctx, cursor, prefix+"*", 0).Result()
	if err != nil {
		logger.Fatalf("Error scanning keys: %s", err)
		return nil, errors.New("cannot scanning keys")
	}
	if len(keys) == 0 {
		return make([]any, 0), nil
	}

	switch op {
	case OpDel:
		_, err = r.Client.Del(ctx, keys...).Result()
	case OpGet:
		vals, err = r.Client.MGet(ctx, keys...).Result()
	}
	if err != nil {
		logger.Fatalf("Error %s keys: %s", op, err)
		return nil, errors.New(fmt.Sprintf("cannot %s", op))
	}
	logger.Printf("val len: %d", len(vals))
	return vals, nil
}

func (r *Redis) CreateChannel(channel string) {
	if r.Channels[channel] != nil {
		err := r.Channels[channel].Close()
		if err != nil {
			logger.Errorf("close exist channel error: %s", err.Error())
		}
	}
	r.Channels[channel] = r.Client.Subscribe(ctx, channel)
}

func (r *Redis) SubscribeChannel(channel string) <-chan *redis.Message {
	ch := r.Channels[channel].Channel()
	return ch
}

func (r *Redis) PublishMessage(channel string, message any) {
	r.Client.Publish(ctx, channel, message)
}

func (r *Redis) CloseChannel(channel string) {
	err := r.Channels[channel].Close()
	if err != nil {
		logger.Errorf("close channel error: %s", err.Error())
		return
	}
}
