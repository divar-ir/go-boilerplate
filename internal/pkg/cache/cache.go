package cache

import (
	context "context"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
)

// PostCache describes a caching mechanism
type PostCache interface {
	// Get item from cache, should return `false` on not found.
	Get(ctx context.Context, token string) (*postview.Post, bool, error)

	// Put item in cache with "at least" `expire` validity time
	Set(ctx context.Context, post *postview.Post, expire time.Duration) error
}
