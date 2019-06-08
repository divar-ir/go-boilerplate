package main

import (
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/metrics"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/metrics/prometheus"
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
