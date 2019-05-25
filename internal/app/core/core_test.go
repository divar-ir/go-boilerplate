package core_test

import (
	"context"
	"errors"
	"testing"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/app/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"

	cacheMocks "git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache/mocks"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	providerMocks "git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CoreTestSuite struct {
	suite.Suite
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}

func (s *CoreTestSuite) TestShouldReturnNotFoundIfProviderReturnsNotFound() {
	cache := &cacheMocks.PostCache{}
	cache.On("Get", mock.Anything, mock.Anything).Once().Return(nil, false, nil)

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
	cache := &cacheMocks.PostCache{}
	cache.On("Get", mock.Anything, "token").Once().Return(&postview.Post{
		Token: "token",
		Title: "title",
	}, true, nil)

	c := core.New(nil, cache)
	response, err := c.GetPost(context.Background(), &postview.GetPostRequest{
		Token: "token",
	})

	s.Nil(err)
	s.Equal(response.Post.Token, "token")
	s.Equal(response.Post.Title, "title")
	cache.AssertExpectations(s.T())
}

func (s *CoreTestSuite) TestShouldReturnFromProviderIfNotCached() {
	cache := &cacheMocks.PostCache{}
	cache.On("Get", mock.Anything, "token").Once().Return(nil, false, nil)
	cache.On("Set", mock.Anything, mock.MatchedBy(func(p *postview.Post) bool {
		s.Equal("token", p.Token)
		s.Equal("title", p.Title)
		return "token" == p.Token && "title" == p.Title
	}), mock.Anything).Once().Return(nil)

	provider := &providerMocks.PostProvider{}
	provider.On("GetPost", mock.Anything, "token").Once().Return(&postview.Post{
		Token: "token",
		Title: "title",
	}, nil)

	c := core.New(provider, cache)
	response, err := c.GetPost(context.Background(), &postview.GetPostRequest{
		Token: "token",
	})

	s.Nil(err)
	s.Equal(response.Post.Token, "token")
	s.Equal(response.Post.Title, "title")
	cache.AssertExpectations(s.T())
	provider.AssertExpectations(s.T())
}

func (s *CoreTestSuite) TestShouldContinueIfCacheFails() {
	cache := &cacheMocks.PostCache{}
	cache.On("Get", mock.Anything, "token").Once().Return(nil, false, errors.New("some error"))
	cache.On("Set", mock.Anything, mock.MatchedBy(func(p *postview.Post) bool {
		s.Equal("token", p.Token)
		s.Equal("title", p.Title)
		return "token" == p.Token && "title" == p.Title
	}), mock.Anything).Once().Return(errors.New("some error"))

	provider := &providerMocks.PostProvider{}
	provider.On("GetPost", mock.Anything, "token").Once().Return(&postview.Post{
		Token: "token",
		Title: "title",
	}, nil)

	c := core.New(provider, cache)
	response, err := c.GetPost(context.Background(), &postview.GetPostRequest{
		Token: "token",
	})

	s.Nil(err)
	s.Equal(response.Post.Token, "token")
	s.Equal(response.Post.Title, "title")
	cache.AssertExpectations(s.T())
	provider.AssertExpectations(s.T())

}
