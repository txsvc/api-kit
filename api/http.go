package api

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/PuerkitoBio/rehttp"

	"github.com/txsvc/apikit/logger"
)

func NewHTTPClient(logger logger.Logger, transport http.RoundTripper) (*http.Client, error) {
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

	loggingTransport := &LoggingTransport{
		InnerTransport: retryTransport,
		Logger:         logger,
	}

	httpClient := &http.Client{
		Transport: loggingTransport,
	}

	return httpClient, nil
}

type LoggingTransport struct {
	InnerTransport http.RoundTripper
	Logger         logger.Logger
}

type contextKey struct {
	name string
}

var contextKeyRequestStart = &contextKey{"RequestStart"}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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

func (t *LoggingTransport) logRequest(req *http.Request) {

	t.Logger.Debugf("--> %s %s\n", req.Method, req.URL)

	if req.Body == nil {
		return
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)

	if err != nil {
		t.Logger.Debug("error reading request body:", err)
	} else {
		t.Logger.Debug(string(data))
	}

	if req.Body != nil {
		t.Logger.Debug(req.Body)
	}

	req.Body = io.NopCloser(bytes.NewReader(data))
}

func (t *LoggingTransport) logResponse(resp *http.Response) {
	ctx := resp.Request.Context()
	defer resp.Body.Close()

	if start, ok := ctx.Value(contextKeyRequestStart).(time.Time); ok {
		t.Logger.Debugf("<-- %d %s (%s)\n", resp.StatusCode, resp.Request.URL, Duration(time.Since(start), 2))
	} else {
		t.Logger.Debugf("<-- %d %s\n", resp.StatusCode, resp.Request.URL)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Logger.Debug("error reading response body:", err)
	} else {
		t.Logger.Debug(string(data))
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
}

func Duration(d time.Duration, dicimal int) time.Duration {
	shift := int(math.Pow10(dicimal))

	units := []time.Duration{time.Second, time.Millisecond, time.Microsecond, time.Nanosecond}
	for _, u := range units {
		if d > u {
			div := u / time.Duration(shift)
			if div == 0 {
				break
			}
			d = d / div * div
			break
		}
	}

	return d
}
