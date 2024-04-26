package storage

import (
	"context"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"minik8s/utils"
	"time"

	"github.com/redis/go-redis/v9"
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

func (r *Redis) RedisSet(key string, value any) error {
	if r.Client == nil {
		r.Client = createRedisClient()
	}
	err := r.Client.Set(ctx, key, value, 500*time.Millisecond)
	if err != nil {
		return errors.New(fmt.Sprintf("set key <%s> failed", key))
	}
	return nil
}

// RedisGet bind should be a pointer
func (r *Redis) RedisGet(key string, bind any) error {
	if r.Client == nil {
		r.Client = createRedisClient()
	}
	val, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		fmt.Println("key2 does not exist")
	}
	if err != nil {
		return err
	}
	err = utils.JsonUnMarshal(val, bind)
	if err != nil {
		return err
	}
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
	var cursor uint64 = 0
	var keys []string
	var vals []any
	for {
		var err error
		keys, cursor, err = r.Client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			logger.Fatalf("Error scanning keys: %s", err)
			return nil, errors.New("cannot scanning keys")
		}
		if len(keys) == 0 {
			break
		}

		switch op {
		case OpDel:
			_, err = r.Client.Del(ctx, keys...).Result()
		case OpGet:
			var gets []any
			gets, err = r.Client.MGet(ctx, keys...).Result()
			vals = append(vals, gets)
		}
		if err != nil {
			logger.Fatalf("Error %s keys: %s", op, err)
			return nil, errors.New(fmt.Sprintf("cannot %s", op))
		}
	}
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
	r.Channels[channel].Close()
}
