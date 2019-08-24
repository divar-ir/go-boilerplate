package main

import (
	"github.com/cafebazaar/go-boilerplate/internal/pkg/metrics"
	"github.com/cafebazaar/go-boilerplate/internal/pkg/metrics/prometheus"
)

var (
	cacheMetrics        metrics.Observer
	postProviderMetrics metrics.Observer
)

func init() {
	cacheMetrics = prometheus.NewHistogram("divar_post_view_cache",
		"view metrics about cache", "cache_type", "method", "ok", "success")

	postProviderMetrics = prometheus.NewHistogram("divar_post_view_provider",
		"view metrics about post provider", "provider_type", "method", "ok", "success")
}

func providePrometheus(config *Config) *prometheus.Server {
	return prometheus.NewServer(config.MetricListenPort)
}
