package core

import (
	"context"

	"github.com/cafebazaar/go-boilerplate/pkg/errors"
	"golang.org/x/xerrors"

	"github.com/cafebazaar/go-boilerplate/internal/app/provider"
	"github.com/cafebazaar/go-boilerplate/pkg/cache"
	"github.com/cafebazaar/go-boilerplate/pkg/postview"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type core struct {
	provider provider.PostProvider
	cache    cache.Layer
}

func New(provider provider.PostProvider, cache cache.Layer) postview.PostViewServer {
	return &core{
		provider: provider,
		cache:    cache,
	}
}

func (c *core) GetPost(ctx context.Context, request *postview.GetPostRequest) (*postview.GetPostResponse, error) {
	post, err := c.getPostFromCache(ctx, request.Token)
	if err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"token": request.Token,
		}).Error("failed to load data from cache")
	} else {
		return &postview.GetPostResponse{
			Post: post,
		}, nil
	}

	post, err = c.provider.GetPost(ctx, request.Token)
	if err != nil {
		if xerrors.Is(err, provider.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "post not found")
		}

		return nil, errors.WrapWithExtra(err, "failed to acquire post", map[string]interface{}{
			"request": request,
		})
	}

	err = c.setPostFromCache(ctx, request.Token, post)
	if err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"token": request.Token,
		}).Error("failed to set data in cache")
	}
	return &postview.GetPostResponse{
		Post: post,
	}, nil
}

func (c *core) getPostFromCache(ctx context.Context, token string) (*postview.Post, error) {
	data, err := c.cache.Get(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get post from cache")
	}
	var result = &postview.Post{}
	err = proto.Unmarshal(data, result)
	if err != nil {
		return nil, errors.WrapWithExtra(err, "failed to unmarshal proto", map[string]interface{}{
			"token": token,
		})
	}

	logrus.Info("load post from cache")
	return result, nil

}

func (c *core) setPostFromCache(ctx context.Context, token string, post *postview.Post) (err error) {
	data, err := proto.Marshal(post)
	if err != nil {
		return errors.Wrap(err, "fail to marshal post")
	}
	err = c.cache.Set(ctx, token, data)
	if err != nil {
		return errors.Wrap(err, "fail to set post in cache")
	}
	return nil
}
