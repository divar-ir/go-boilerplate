package provider

import (
	context "context"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
)

// PostProvider specifies mechanism of retrieving posts.
// which can be either a DB, some microservice, etc.
type PostProvider interface {
	// Get post with requested token
	GetPost(ctx context.Context, token string) (*postview.Post, error)

	// Add post to datastore
	AddPost(ctx context.Context, post *postview.Post) error
}
