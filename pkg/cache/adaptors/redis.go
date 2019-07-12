package adaptors

import (
	"context"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

type redisCache struct {
	ins        *redis.Client
	expireTime time.Duration
}

func NewRedisAdaptor(expireTime time.Duration, redisClient *redis.Client) cache.Layer {
	return &redisCache{
		expireTime: expireTime,
		ins:        redisClient,
	}
}

func (rc *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if rc == nil {
		return nil, errors.New("redis cache is disabled")
	}
	client := rc.ins.WithContext(ctx)
	value, err := client.Get(key).Result()
	return []byte(value), err
}

func (rc *redisCache) Set(ctx context.Context, key string, value []byte) error {
	if rc == nil {
		return errors.New("redis cache is disabled")
	}
	client := rc.ins.WithContext(ctx)
	err := client.SetNX(key, value, rc.expireTime).Err()
	return err
}

func (rc *redisCache) Delete(ctx context.Context, key string) error {
	if rc == nil {
		return errors.New("redis cache is disable")
	}
	client := rc.ins.WithContext(ctx)
	err := client.Del(key).Err()
	return err
}

func (rc *redisCache) Clear(ctx context.Context) error {
	if rc == nil {
		return errors.New("redis client is disable")
	}
	client := rc.ins.WithContext(ctx)
	err := client.FlushDB().Err()
	return err
}
