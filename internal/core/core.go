package core

import (
	"context"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type core struct {
	provider postview.PostProvider
	cache    postview.PostCache
}

func New(provider postview.PostProvider, cache postview.PostCache) postview.PostViewServer {
	return &core{
		provider: provider,
		cache:    cache,
	}
}

func (c *core) GetPost(ctx context.Context, request *postview.GetPostRequest) (*postview.GetPostResponse, error) {
	post, err := c.provider.GetPost(ctx, request.Token)
	if err != nil {
		if err == postview.ErrNotFound {
			return nil, status.Error(codes.NotFound, "post not found")
		}

		return nil, err
	}

	return &postview.GetPostResponse{
		Error: 0,
		Post:  post,
	}, nil
}
