package provider

import (
	"context"
	"sync"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
)

type memoryProvider struct {
	items sync.Map
}

func NewMemory() postview.PostProvider {
	return &memoryProvider{}
}

func (provider *memoryProvider) GetPost(ctx context.Context, token string) (*postview.Post, error) {
	value, ok := provider.items.Load(token)
	if !ok {
		return nil, postview.ErrNotFound
	}

	return value.(*postview.Post), nil
}

func (provider *memoryProvider) Add(post *postview.Post) {
	provider.items.Store(post.Token, post)
}
