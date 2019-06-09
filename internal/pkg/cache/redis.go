package cache

import (
	context "context"
	"fmt"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/errors"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

type redisCache struct {
	client *redis.Client
	prefix string
}

func NewRedis(client *redis.Client, prefix string) PostCache {
	return redisCache{
		client: client,
		prefix: prefix,
	}
}

func (c redisCache) Get(ctx context.Context, token string) (*postview.Post, bool, error) {
	key := c.generateKey(token)
	data, err := c.client.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, errors.WrapWithExtra(err, "failed to read from redis", map[string]interface{}{
			"token": token,
		})
	}

	var result postview.Post
	err = proto.Unmarshal(data, &result)
	if err != nil {
		return nil, false, errors.WrapWithExtra(err, "failed to unmarshal proto", map[string]interface{}{
			"token": token,
		})
	}

	return &result, true, nil
}

func (c redisCache) Set(ctx context.Context, post *postview.Post, expire time.Duration) error {
	data, err := proto.Marshal(post)
	if err != nil {
		return errors.WrapWithExtra(err, "failed to marshal post", map[string]interface{}{
			"post": post,
		})
	}

	key := c.generateKey(post.Token)

	err = c.client.Set(key, data, expire).Err()
	if err != nil {
		return errors.WrapWithExtra(err, "failed to write to redis", map[string]interface{}{
			"post": post,
		})
	}
	return nil
}

func (c redisCache) generateKey(token string) string {
	return fmt.Sprintf("%s:%s", c.prefix, token)
}
