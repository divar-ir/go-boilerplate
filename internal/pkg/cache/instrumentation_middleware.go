package cache

import (
	context "context"
	"fmt"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/metrics"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
)

type instrumentationMiddleware struct {
	next   PostCache
	timing metrics.Observer
}

func NewInstrumentationMiddleware(next PostCache, timing metrics.Observer) PostCache {
	return instrumentationMiddleware{
		next:   next,
		timing: timing,
	}
}

func (m instrumentationMiddleware) Get(
	ctx context.Context, token string) (result *postview.Post, ok bool, err error) {

	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"ok":      fmt.Sprint(ok),
			"success": fmt.Sprint(err == nil),
			"method":  "Get",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.Get(ctx, token)
}

func (m instrumentationMiddleware) Set(
	ctx context.Context, post *postview.Post, expire time.Duration) (err error) {

	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "Set",
			"ok":      "true",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.Set(ctx, post, expire)
}
