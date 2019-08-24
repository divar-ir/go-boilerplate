package prometheus

import (
	"github.com/cafebazaar/go-boilerplate/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type histogram struct {
	vec         *prometheus.HistogramVec
	labels      []string
	labelValues map[string]string
}

func NewHistogram(name, help string, labels ...string) metrics.Observer {
	result := &histogram{
		vec: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: name,
			Help: help,
		}, labels),

		labels:      labels,
		labelValues: make(map[string]string),
	}

	prometheus.MustRegister(result.vec)

	return result
}

func (s *histogram) Observe(value float64) {
	values := make([]string, len(s.labels))

	for i, label := range s.labels {
		values[i] = s.labelValues[label]
	}

	s.vec.WithLabelValues(values...).Observe(value)
}

func (s *histogram) With(labels map[string]string) metrics.Observer {
	newLabelValues := make(map[string]string)

	for label, value := range s.labelValues {
		newLabelValues[label] = value
	}

	for label, value := range labels {
		newLabelValues[label] = value
	}

	return &histogram{
		labelValues: newLabelValues,
		labels:      s.labels,
		vec:         s.vec,
	}
}
