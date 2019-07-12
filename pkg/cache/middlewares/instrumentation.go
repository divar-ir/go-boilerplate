package middlewares

import (
	"context"
	"fmt"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/metrics"
)

type instrumentationMiddleware struct {
	next   cache.Layer
	timing metrics.Observer
}

func NewInstrumentationMiddleware(layer cache.Layer, timing metrics.Observer) cache.Layer {
	return instrumentationMiddleware{
		next:   layer,
		timing: timing,
	}
}

func (m instrumentationMiddleware) Get(ctx context.Context, key string) (data []byte, err error) {
	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "Get",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.Get(ctx, key)
}

func (m instrumentationMiddleware) Set(ctx context.Context, key string, value []byte) (err error) {
	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "Set",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.Set(ctx, key, value)
}

func (m instrumentationMiddleware) Delete(ctx context.Context, key string) (err error) {
	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "Delete",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.Delete(ctx, key)
}

func (m instrumentationMiddleware) Clear(ctx context.Context) (err error) {
	return m.next.Clear(ctx)
}
