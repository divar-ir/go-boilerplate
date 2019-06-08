package prometheus_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func waitForMetric(t *testing.T, metric string) {
	s := assert.New(t)

	hasAppeared := func() bool {
		return len(filterOnMetric(t, metric)) > 0
	}

	if hasAppeared() {
		return
	}

	timeout := time.NewTicker(5 * time.Second)
	defer timeout.Stop()

	interval := time.NewTicker(100 * time.Millisecond)
	defer interval.Stop()

	for {
		select {
		case <-timeout.C:
			s.FailNowf("metric did not appear",
				"metric %s did not appear on prometheus", metric)
			return

		case <-interval.C:
			if hasAppeared() {
				return
			}
		}
	}
}

func filterOnMetric(t *testing.T, metric string) []string {
	var result []string

	for _, line := range strings.Split(dumpPrometheus(t), "\n") {
		if strings.Contains(line, metric) {
			result = append(result, line)
		}
	}

	return result
}

func dumpPrometheus(t *testing.T) string {
	s := assert.New(t)

	responseWriter := newFakeResponseWriter()
	request := &http.Request{
		Method:     "GET",
		RequestURI: "/",
	}

	promhttp.Handler().ServeHTTP(responseWriter, request)

	if responseWriter.statusCode != 200 {
		s.FailNowf("prometheus handler failed",
			"expected status code 200 but got %d", responseWriter.statusCode)

		return ""
	}

	return responseWriter.body.String()
}

type fakeResponseWriter struct {
	statusCode int
	body       bytes.Buffer
	headers    http.Header
}

func newFakeResponseWriter() *fakeResponseWriter {
	return &fakeResponseWriter{
		headers: make(http.Header),
	}
}

func (f *fakeResponseWriter) Header() http.Header {
	return f.headers
}

func (f *fakeResponseWriter) Write(body []byte) (int, error) {
	return f.body.Write(body)
}

func (f *fakeResponseWriter) WriteHeader(statusCode int) {
	f.statusCode = statusCode
}
