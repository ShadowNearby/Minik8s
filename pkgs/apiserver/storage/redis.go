package storage

import (
	"context"
	"errors"
	"fmt"
	"minik8s/config"
	"minik8s/pkgs/constants"
	"minik8s/utils"

	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"
)

type Redis struct {
	Client   *redis.Client
	Channels map[string]*redis.PubSub
}

var ctx = context.Background()

const (
	OpDel string = "delete"
	OpGet string = "get"
)

func CreateRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", config.ClusterMasterIP),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
func createFunctionClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:8070",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (r *Redis) InitChannels() {
	for _, channel := range constants.Channels {
		for _, operation := range constants.Operations {
			name := constants.GenerateChannelName(channel, operation)
			r.CreateChannel(name)
			logger.Infof("channel %s created", name)
		}
	}
	for _, channel := range constants.OtherChannels {
		r.CreateChannel(channel)
	}
}

func (r *Redis) redisSet(key string, value any) error {
	if _, ok := value.(string); !ok {
		value = utils.JsonMarshal(value)
	}
	err := r.Client.Set(ctx, key, value, 0)
	if err != nil {
		return err.Err()
	}
	return nil
}

// redisGet bind should be a pointer, only for json object
func (r *Redis) redisGet(key string, ptr any) error {
	val, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		logger.Errorf("%s does not exist", key)
	}
	if err != nil {
		return err
	}
	err = utils.JsonUnMarshal(val, ptr)
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
	logger.Debugf("in redis range op")
	var cursor uint64 = 0
	var keys []string
	var vals []any
	var err error
	keys, _, err = r.Client.Scan(ctx, cursor, prefix+"*", 0).Result()
	if err != nil {
		logger.Errorf("Error scanning keys: %s", err)
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
		logger.Errorf("Error %s keys: %s", op, err)
		return nil, fmt.Errorf("cannot %s", op)
	}
	logger.Debugf("val len: %d", len(vals))
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
	logger.Infoln("Subscribe channel:", channel)
	ch := r.Channels[channel].Channel()
	return ch
}

func (r *Redis) PublishMessage(channel string, message any) {
	// logger.Infoln("[", channel, "]\t", message)
	r.Client.Publish(ctx, channel, message)
}

func (r *Redis) CloseChannel(channel string) {
	err := r.Channels[channel].Close()
	if err != nil {
		logger.Errorf("close channel error: %s", err.Error())
		return
	}
}
