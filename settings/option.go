// Most of the code is lifted from
// https://github.com/googleapis/google-api-go-client/blob/main/option/option.go
//
// For details and copyright etc. see above url.
//
package settings

type (
	// ClientOption implements a configuration option for an API client.
	ClientOption interface {
		Apply(*DialSettings)
	}
)

// WithEndpoint returns a ClientOption that overrides the default endpoint to be used for a service.
func WithEndpoint(url string) ClientOption {
	return withEndpoint(url)
}

type withEndpoint string

func (w withEndpoint) Apply(o *DialSettings) {
	o.Endpoint = string(w)
}

// WithScopes returns a ClientOption that overrides the default (OAuth2) scopes to be used for a service.
func WithScopes(scope ...string) ClientOption {
	return withScopes(scope)
}

type withScopes []string

func (w withScopes) Apply(o *DialSettings) {
	o.Scopes = make([]string, len(w))
	copy(o.Scopes, w)
}

// WithUserAgent returns a ClientOption that sets the User-Agent.
func WithUserAgent(ua string) ClientOption {
	return withUA(ua)
}

type withUA string

func (w withUA) Apply(o *DialSettings) { o.UserAgent = string(w) }

// WithAPIKey returns a ClientOption that specifies an API key to be used as the basis for authentication.
func WithAPIKey(apiKey string) ClientOption {
	return withAPIKey(apiKey)
}

type withAPIKey string

func (w withAPIKey) Apply(o *DialSettings) { o.APIKey = string(w) }

// WithoutAuthentication returns a ClientOption that specifies that no
// authentication should be used. It is suitable only for testing and for
// accessing public resources.
// It is an error to provide both WithoutAuthentication and any of WithAPIKey,
// WithTokenSource, WithCredentialsFile or WithServiceAccountFile.
func WithoutAuthentication() ClientOption {
	return withoutAuthentication{}
}

type withoutAuthentication struct{}

func (w withoutAuthentication) Apply(o *DialSettings) { o.NoAuth = true }
