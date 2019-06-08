package cache_test

import (
	context "context"
	"errors"
	"testing"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache"
	cacheMocks "git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache/mocks"
	metricsMocks "git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/metrics/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CacheInstrumentationMiddlewareTestSuite struct {
	suite.Suite
}

func TestCacheInstrumentationMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CacheInstrumentationMiddlewareTestSuite))
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldCountOkIsFalse() {
	mockedCache := &cacheMocks.PostCache{}
	mockedCache.On("Get", mock.Anything, mock.Anything).Return(nil, false, nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method": "Get",
		"ok":     "false",
	})

	hooked := cache.NewInstrumentationMiddleware(mockedCache, mockedObserver)
	_, ok, err := hooked.Get(context.Background(), "")
	s.False(ok)
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedCache.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldCountOkIsTrue() {
	mockedCache := &cacheMocks.PostCache{}
	mockedCache.On("Get", mock.Anything, mock.Anything).Return(nil, true, nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method": "Get",
		"ok":     "true",
	})

	hooked := cache.NewInstrumentationMiddleware(mockedCache, mockedObserver)
	_, ok, err := hooked.Get(context.Background(), "")
	s.True(ok)
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedCache.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldErrIsNil() {
	mockedCache := &cacheMocks.PostCache{}
	mockedCache.On("Get", mock.Anything, mock.Anything).Return(nil, false, nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "Get",
		"success": "true",
	})

	hooked := cache.NewInstrumentationMiddleware(mockedCache, mockedObserver)
	_, _, err := hooked.Get(context.Background(), "")
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedCache.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldErrIsNotNil() {
	mockedCache := &cacheMocks.PostCache{}
	mockedCache.On("Get", mock.Anything, mock.Anything).Return(nil, false, errors.New("some err"))

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "Get",
		"success": "false",
	})

	hooked := cache.NewInstrumentationMiddleware(mockedCache, mockedObserver)
	_, _, err := hooked.Get(context.Background(), "")
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedCache.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestSetShouldCountErrIsNil() {
	mockedCache := &cacheMocks.PostCache{}
	mockedCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "Set",
		"success": "true",
	})

	hooked := cache.NewInstrumentationMiddleware(mockedCache, mockedObserver)
	err := hooked.Set(context.Background(), nil, 0)
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedCache.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestSetShouldCountErrIsNotNil() {
	mockedCache := &cacheMocks.PostCache{}
	mockedCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some err"))

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "Set",
		"success": "false",
	})

	hooked := cache.NewInstrumentationMiddleware(mockedCache, mockedObserver)
	err := hooked.Set(context.Background(), nil, 0)
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedCache.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) makeObserver(
	expectedLabels map[string]string) *metricsMocks.Observer {

	mockedObserver := &metricsMocks.Observer{}
	mockedObserver.On("Observe", mock.Anything).Once()
	mockedObserver.On("With", s.makeMatcher(expectedLabels)).Once().Return(mockedObserver)

	return mockedObserver
}

func (s *CacheInstrumentationMiddlewareTestSuite) makeMatcher(
	expectedLabels map[string]string) interface{} {

	return mock.MatchedBy(func(labels map[string]string) bool {
		result := true

		for expectedKey, expectedValue := range expectedLabels {
			value, ok := labels[expectedKey]
			s.True(ok, "expected to find label %v", expectedKey)
			if !ok {
				result = false
				continue
			}

			s.Equal(expectedValue, value,
				"expected to find value %v for key %v", expectedValue, expectedKey)
			result = result && expectedValue == value
		}

		return result
	})
}
