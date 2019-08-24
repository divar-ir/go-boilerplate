package core_test

import (
	"context"
	"github.com/cafebazaar/go-boilerplate/internal/pkg/provider"
	"github.com/cafebazaar/go-boilerplate/pkg/cache/adaptors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/cafebazaar/go-boilerplate/internal/app/core"
	"github.com/cafebazaar/go-boilerplate/pkg/postview"
	"github.com/golang/protobuf/proto"

	providerMocks "github.com/cafebazaar/go-boilerplate/internal/pkg/provider/mocks"
	"github.com/stretchr/testify/suite"
)

type CoreTestSuite struct {
	suite.Suite
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}

func (s *CoreTestSuite) TestShouldReturnNotFoundIfProviderReturnsNotFound() {
	cache := adaptors.NewSynMapAdaptor(logrus.New())

	mockProvider := &providerMocks.PostProvider{}
	mockProvider.On("GetPost", mock.Anything, mock.Anything).Once().Return(nil, provider.ErrNotFound)

	c := core.New(mockProvider, cache)
	_, err := c.GetPost(context.Background(), &postview.GetPostRequest{})
	s.NotNil(err)

	grpcStatus, ok := status.FromError(err)
	s.True(ok)
	if !ok {
		return
	}
	s.Equal(codes.NotFound, grpcStatus.Code())
}

func (s *CoreTestSuite) TestShouldReturnFromCacheIfFound() {
	cache := adaptors.NewSynMapAdaptor(logrus.New())
	c := core.New(nil, cache)
	token := "token"
	title := "title"
	data, err := proto.Marshal(&postview.Post{
		Token: token,
		Title: title,
	})
	if !s.NoError(err, "fail to marshalize post") {
		return
	}
	err = cache.Set(context.Background(), token, data)
	if !s.NoError(err, "fail to marshalize post") {
		return
	}
	response, err := c.GetPost(context.Background(), &postview.GetPostRequest{
		Token: token,
	})

	s.NoError(err)
	s.Equal(response.Post.Token, token)
	s.Equal(response.Post.Title, title)
}

func (s *CoreTestSuite) TestShouldReturnFromProviderIfNotCached() {
	cache := adaptors.NewSynMapAdaptor(logrus.New())

	mockProvider := &providerMocks.PostProvider{}
	mockProvider.On("GetPost", mock.Anything, "token").Once().Return(&postview.Post{
		Token: "token",
		Title: "title",
	}, nil)

	c := core.New(mockProvider, cache)
	response, err := c.GetPost(context.Background(), &postview.GetPostRequest{
		Token: "token",
	})

	s.Nil(err)
	s.Equal(response.Post.Token, "token")
	s.Equal(response.Post.Title, "title")
}
