package provider

import (
	"context"
	"sync"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/errors"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
)

type memoryProvider struct {
	items sync.Map
}

func NewMemory() PostProvider {
	return &memoryProvider{}
}

func (provider *memoryProvider) GetPost(ctx context.Context, token string) (*postview.Post, error) {
	value, ok := provider.items.Load(token)
	if !ok {
		return nil, errors.WrapWithExtra(ErrNotFound, "post not found", map[string]interface{}{
			"token": token,
		})
	}

	return value.(*postview.Post), nil
}

func (provider *memoryProvider) AddPost(ctx context.Context, post *postview.Post) error {
	provider.items.Store(post.Token, post)
	return nil
}
