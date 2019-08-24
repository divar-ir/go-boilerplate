package middlewares

import (
	"context"
	"github.com/cafebazaar/go-boilerplate/pkg/cache/adaptors"
	"github.com/sirupsen/logrus"
	"testing"

	metricsMocks "github.com/cafebazaar/go-boilerplate/internal/pkg/metrics/mocks"

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
	cache := adaptors.NewSynMapAdaptor(logrus.New())

	mockedObserver := s.makeObserver(map[string]string{
		"method": "Get",
	})

	hooked := NewInstrumentationMiddleware(cache, mockedObserver)
	_, err := hooked.Get(context.Background(), "")
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldCountOkIsTrue() {
	cache := adaptors.NewSynMapAdaptor(logrus.New())
	err := cache.Set(context.Background(), "", []byte("value"))
	s.NoError(err)

	mockedObserver := s.makeObserver(map[string]string{
		"method": "Get",
	})

	hooked := NewInstrumentationMiddleware(cache, mockedObserver)
	_, err = hooked.Get(context.Background(), "")
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
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
