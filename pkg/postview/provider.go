package postview

import context "context"

// PostProvider specifies mechanism of retrieving posts.
// which can be either a DB, some microservice, etc.
type PostProvider interface {
	// Get post with requested token
	GetPost(ctx context.Context, token string) (*Post, error)
}
