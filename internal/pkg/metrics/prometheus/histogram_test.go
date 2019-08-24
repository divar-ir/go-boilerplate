package prometheus_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/cafebazaar/go-boilerplate/internal/pkg/metrics"
	"github.com/cafebazaar/go-boilerplate/internal/pkg/metrics/prometheus"
	"github.com/stretchr/testify/suite"
)

type HistogramTestSuite struct {
	suite.Suite

	sampleMetric          metrics.Observer
	sampleMetricWithLabel metrics.Observer
}

func TestHistogramTestSuite(t *testing.T) {
	suite.Run(t, new(HistogramTestSuite))
}

func (s *HistogramTestSuite) TestShouldContainType() {
	s.True(strings.Contains(dumpPrometheus(s.T()), "# TYPE sample_histogram histogram"))
}

func (s *HistogramTestSuite) TestShouldContainHelp() {
	s.True(strings.Contains(dumpPrometheus(s.T()), "# HELP sample_histogram a sample metric"))
}

func (s *HistogramTestSuite) TestShouldCount() {
	s.True(strings.Contains(dumpPrometheus(s.T()), "sample_histogram_count 2"))
}

func (s *HistogramTestSuite) TestShouldSum() {
	s.True(strings.Contains(dumpPrometheus(s.T()), "sample_histogram_sum 3"))
}

func (s *HistogramTestSuite) TestShouldApplyLabels() {
	s.True(strings.Contains(dumpPrometheus(s.T()), "sample_histogram_with_label_count{my_label=\"my_value\"} 1"))
}

func (s *HistogramTestSuite) TestShouldNotHaveDataRaceOnConcurrentAccessWithLabel() {
	const NumberOfGoRoutines = 10

	metric := prometheus.NewHistogram("concurrent_histogram_metric_with_label", "a sample metric", "my_label")
	var started sync.WaitGroup
	var beginOperating sync.WaitGroup
	var testFinished sync.WaitGroup

	started.Add(NumberOfGoRoutines)
	testFinished.Add(NumberOfGoRoutines)
	beginOperating.Add(1)

	for i := 0; i < NumberOfGoRoutines; i++ {
		go func() {
			defer testFinished.Done()

			started.Done()
			beginOperating.Wait()

			for j := 0; j < 1000; j++ {
				metric.With(map[string]string{
					"my_label": "my_value",
				}).Observe(1.0)
			}
		}()
	}

	started.Wait()
	beginOperating.Done()

	testFinished.Wait()
}

func (s *HistogramTestSuite) SetupSuite() {
	s.sampleMetricWithLabel = prometheus.NewHistogram("sample_histogram_with_label", "a sample metric", "my_label")
	s.sampleMetric = prometheus.NewHistogram("sample_histogram", "a sample metric")

	s.sampleMetric.Observe(1.0)
	s.sampleMetric.Observe(2.0)
	waitForMetric(s.T(), "sample_histogram")

	s.sampleMetricWithLabel.With(map[string]string{
		"my_label": "my_value",
	}).Observe(1.0)
	waitForMetric(s.T(), "sample_histogram_with_label")
}
