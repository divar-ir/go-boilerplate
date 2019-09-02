package provider_test

import (
	"context"
	"errors"
	"github.com/cafebazaar/go-boilerplate/internal/app/provider"
	"testing"

	providerMocks "github.com/cafebazaar/go-boilerplate/internal/app/provider/mocks"
	metricsMocks "github.com/cafebazaar/go-boilerplate/internal/pkg/metrics/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProviderInstrumentationMiddlewareTestSuite struct {
	suite.Suite
}

func TestProviderInstrumentationMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderInstrumentationMiddlewareTestSuite))
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestGetPostShouldCountErrIsNil() {
	mockedProvider := &providerMocks.PostProvider{}
	mockedProvider.On("GetPost", mock.Anything, mock.Anything).Return(nil, nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "GetPost",
		"success": "true",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	_, err := hooked.GetPost(context.Background(), "")
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestGetPostShouldCountErrIsNotNil() {
	mockedProvider := &providerMocks.PostProvider{}
	mockedProvider.On("GetPost", mock.Anything, mock.Anything).Return(nil, errors.New("some err"))

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "GetPost",
		"success": "false",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	_, err := hooked.GetPost(context.Background(), "")
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestAddPostShouldCountErrIsNil() {
	mockedProvider := &providerMocks.PostProvider{}
	mockedProvider.On("AddPost", mock.Anything, mock.Anything).Return(nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "AddPost",
		"success": "true",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	err := hooked.AddPost(context.Background(), nil)
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestAddPostShouldCountErrIsNotNil() {
	mockedProvider := &providerMocks.PostProvider{}
	mockedProvider.On("AddPost", mock.Anything, mock.Anything).Return(errors.New("some err"))

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "AddPost",
		"success": "false",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	err := hooked.AddPost(context.Background(), nil)
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) makeObserver(
	expectedLabels map[string]string) *metricsMocks.Observer {

	mockedObserver := &metricsMocks.Observer{}
	mockedObserver.On("Observe", mock.Anything).Once()
	mockedObserver.On("With", s.makeMatcher(expectedLabels)).Once().Return(mockedObserver)

	return mockedObserver
}

func (s *ProviderInstrumentationMiddlewareTestSuite) makeMatcher(
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
