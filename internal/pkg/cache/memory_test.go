package cache_test

import (
	context "context"
	"testing"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache"

	"github.com/stretchr/testify/suite"
)

type MemoryCacheTestSuite struct {
	suite.Suite

	cache cache.PostCache
}

func TestMemoryCacheTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryCacheTestSuite))
}

func (s *MemoryCacheTestSuite) TestShouldReturnFalseInitially() {
	_, ok, err := s.cache.Get(context.Background(), "token")
	s.Nil(err)
	s.False(ok)
}

func (s *MemoryCacheTestSuite) TestShouldReturnTrueAfterSet() {
	s.cache.Set(context.Background(), &postview.Post{
		Title: "a title",
		Token: "token",
	}, 1*time.Second)
	_, ok, err := s.cache.Get(context.Background(), "token")
	s.Nil(err)
	s.True(ok)
}

func (s *MemoryCacheTestSuite) TestShouldGetAfterSet() {
	s.cache.Set(context.Background(), &postview.Post{
		Title: "a title",
		Token: "token",
	}, 1*time.Second)
	post, _, err := s.cache.Get(context.Background(), "token")
	s.Nil(err)
	s.Equal("a title", post.Title)
	s.Equal("token", post.Token)
}

func (s *MemoryCacheTestSuite) SetupTest() {
	s.cache = cache.NewMemory()
}
