package api

import (
	"net/http"
	"strconv"
	"strings"
)

type (
	// RelevantHeaders represents the most important headers
	RelevantHeaders struct {
		Range           string `header:"Range"`
		UserAgent       string `header:"User-Agent"`
		Forwarded       string `header:"Forwarded"`
		XForwardedFor   string `header:"X-Forwarded-For"`
		XForwwardedHost string `header:"X-Forwarded-Host"`
		Referer         string `header:"Referer"`
	}
)

func (h *RelevantHeaders) Ranges() (int64, int64) {
	return ParseRange(h.Range)
}

// ExtractHeaders extracts the relevant HTTP header stuff only
func ExtractHeaders(r *http.Request) RelevantHeaders {
	h := RelevantHeaders{
		Range:           r.Header.Get("Range"),
		UserAgent:       r.Header.Get("User-Agent"),
		Forwarded:       r.Header.Get("Forwarded"),
		XForwardedFor:   r.Header.Get("X-Forwarded-For"),
		XForwwardedHost: r.Header.Get("X-Forwarded-Host"),
		Referer:         r.Header.Get("Referer"),
	}
	return h
}

// ParseRange extracts a byte range if specified. For specs see
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
func ParseRange(r string) (int64, int64) {
	if r == "" {
		return 0, -1 // no range requested
	}
	parts := strings.Split(r, "=")
	if len(parts) != 2 {
		return 0, -1 // no range requested
	}
	// we simply assume that parts[0] == "bytes"
	ra := strings.Split(parts[1], "-")
	if len(ra) != 2 { // again a simplification, multiple ranges or overlapping ranges are not supported
		return 0, -1
	}

	start, err := strconv.ParseInt(ra[0], 10, 64)
	if err != nil {
		return 0, -1
	}
	end, err := strconv.ParseInt(ra[1], 10, 64)
	if err != nil {
		return 0, -1
	}

	return start, end - start
}
