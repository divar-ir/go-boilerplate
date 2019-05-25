package provider_test

import (
	"context"
	"testing"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"github.com/stretchr/testify/suite"
)

type MemoryProviderTestSuite struct {
	suite.Suite

	provider provider.PostProvider
}

func TestMemoryProviderTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryProviderTestSuite))
}

func (s *MemoryProviderTestSuite) TestShouldReturnNotFoundInitially() {
	_, err := s.provider.GetPost(context.Background(), "token")
	s.Equal(provider.ErrNotFound, err)
}

func (s *MemoryProviderTestSuite) TestShouldReturnPostAfterAdd() {
	err := s.provider.AddPost(context.Background(), &postview.Post{
		Token: "abcd",
		Title: "a title",
	})
	s.Nil(err)
	if err != nil {
		return
	}

	post, err := s.provider.GetPost(context.Background(), "abcd")
	s.Nil(err)
	if err != nil {
		return
	}

	s.Equal("abcd", post.Token)
	s.Equal("a title", post.Title)
}

func (s *MemoryProviderTestSuite) SetupTest() {
	s.provider = provider.NewMemory()
}
