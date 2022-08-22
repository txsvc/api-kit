package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	opts := []ClientOption{
		WithEndpoint("https://example.com:443"),
		WithScopes("a"), // the next WithScopes should overwrite this one
		WithScopes("https://example.com/auth/helloworld", "https://example.com/auth/otherthing"),
		WithUserAgent("ua"),
		WithAPIKey("api-key"),
	}

	var got DialSettings
	for _, opt := range opts {
		opt.Apply(&got)
	}

	want := DialSettings{
		Endpoint:  "https://example.com:443",
		Scopes:    []string{"https://example.com/auth/helloworld", "https://example.com/auth/otherthing"},
		UserAgent: "ua",
		APIKey:    "api-key",
	}

	assert.Equal(t, want, got)
}
