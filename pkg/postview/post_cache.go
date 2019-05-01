package postview

import (
	context "context"
	"time"
)

// PostCache describes a caching mechanism
type PostCache interface {
	// Get item from cache, should return `false` on not found.
	Get(ctx context.Context, token string) (*Post, bool, error)

	// Put item in cache with "at least" `expire` validity time
	Set(ctx context.Context, post *Post, expire time.Duration) error
}
