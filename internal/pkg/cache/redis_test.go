package cache_test

import (
	context "context"
	"strings"
	"testing"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/suite"
)

type RedisCacheTestSuite struct {
	suite.Suite

	db    *miniredis.Miniredis
	cache cache.PostCache
}

func TestRedisCacheTestSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheTestSuite))
}

func (s *RedisCacheTestSuite) TestShouldReturnFalseInitially() {
	_, ok, err := s.cache.Get(context.Background(), "token")
	s.Nil(err)
	s.False(ok)
}

func (s *RedisCacheTestSuite) TestShouldReturnTrueAfterSet() {
	s.Nil(s.cache.Set(context.Background(), &postview.Post{
		Title: "a title",
		Token: "token",
	}, 1*time.Second))
	_, ok, err := s.cache.Get(context.Background(), "token")
	s.Nil(err)
	s.True(ok)
}

func (s *RedisCacheTestSuite) TestShouldGetAfterSet() {
	s.Nil(s.cache.Set(context.Background(), &postview.Post{
		Title: "a title",
		Token: "token",
	}, 1*time.Second))
	post, _, err := s.cache.Get(context.Background(), "token")
	s.Nil(err)
	s.Equal("a title", post.Title)
	s.Equal("token", post.Token)
}

func (s *RedisCacheTestSuite) SetupTest() {
	var err error
	currentTest := strings.Replace(s.T().Name(), "/", "_", -1)

	s.db, err = miniredis.Run()
	if err != nil {
		s.FailNow("failed to create miniredis db")
		return
	}

	client := redis.NewClient(&redis.Options{Addr: s.db.Addr()})
	s.cache = cache.NewRedis(client, currentTest)
}
