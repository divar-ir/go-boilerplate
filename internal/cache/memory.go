package cache

import (
	"context"
	"sync"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
)

type memoryCache struct {
	items sync.Map
}

func NewMemory() postview.PostCache {
	return &memoryCache{}
}

func (provider *memoryCache) Get(ctx context.Context, token string) (*postview.Post, bool, error) {
	value, ok := provider.items.Load(token)
	if !ok {
		return nil, false, nil
	}

	return value.(*postview.Post), true, nil
}

func (provider *memoryCache) Set(ctx context.Context, post *postview.Post, expire time.Duration) error {
	provider.items.Store(post.Token, post)
	return nil
}
