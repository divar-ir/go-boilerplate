package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/cafebazaar/go-boilerplate/internal/pkg/metrics"
	"github.com/cafebazaar/go-boilerplate/pkg/postview"
)

type instrumentationMiddleware struct {
	next   PostProvider
	timing metrics.Observer
}

func NewInstrumentationMiddleware(
	next PostProvider, timing metrics.Observer) PostProvider {

	return instrumentationMiddleware{
		next:   next,
		timing: timing,
	}
}

func (m instrumentationMiddleware) GetPost(
	ctx context.Context, token string) (result *postview.Post, err error) {

	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "GetPost",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.GetPost(ctx, token)
}

func (m instrumentationMiddleware) AddPost(
	ctx context.Context, post *postview.Post) (err error) {

	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "AddPost",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.AddPost(ctx, post)
}
