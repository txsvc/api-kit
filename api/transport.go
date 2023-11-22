package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/txsvc/cloudlib/observer"
)

type (
	loggingTransport struct {
		InnerTransport http.RoundTripper
	}

	contextKey struct {
		name string
	}
)

var contextKeyRequestStart = &contextKey{"RequestStart"}

func NewTransport(transport http.RoundTripper) *http.Client {
	retryTransport := rehttp.NewTransport(
		transport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(3),
			rehttp.RetryAny(
				rehttp.RetryTemporaryErr(),
				rehttp.RetryStatuses(502, 503),
			),
		),
		rehttp.ExpJitterDelay(100*time.Millisecond, 1*time.Second),
	)

	return &http.Client{
		Transport: &loggingTransport{
			InnerTransport: retryTransport,
		},
	}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), contextKeyRequestStart, time.Now())
	req = req.WithContext(ctx)

	t.logRequest(req)

	resp, err := t.InnerTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *loggingTransport) logRequest(req *http.Request) {

	observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("--> %s %s\n", req.Method, req.URL))

	if req.Body == nil {
		return
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)

	if err != nil {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("error reading request body: %v", err))
	} else {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("%v", data))
	}

	if req.Body != nil {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("%v", req.Body))
	}

	req.Body = io.NopCloser(bytes.NewReader(data))
}

func (t *loggingTransport) logResponse(resp *http.Response) {
	ctx := resp.Request.Context()
	defer resp.Body.Close()

	if start, ok := ctx.Value(contextKeyRequestStart).(time.Time); ok {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("<-- %d %s (%s)\n", resp.StatusCode, resp.Request.URL, Duration(time.Since(start), 2)))
	} else {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("<-- %d %s\n", resp.StatusCode, resp.Request.URL))
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("error reading response body: %v", err))
	} else {
		observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("%v", data))
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
}
