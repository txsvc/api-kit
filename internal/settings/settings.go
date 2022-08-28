// Most of the code is lifted from
// https://github.com/googleapis/google-api-go-client/blob/main/internal/settings.go
//
// For details and copyright etc. see above url.
package settings

import (
	"io/fs"
)

const (
	StateInit State = iota - 2
	StateInvalid
	StateUndefined  // logged out
	StateAuthorized // logged in

	indentChar             = "  "
	filePerm   fs.FileMode = 0644
)

type (
	State int

	// Settings holds information needed to establish a connection with a
	// backend API service or to simply configure a service or CLI.
	Settings struct {
		Endpoint string `json:"endpoint,omitempty"`

		Credentials *Credentials `json:"credentials,omitempty"`

		Scopes        []string `json:"scopes,omitempty"`
		DefaultScopes []string `json:"default_scopes,omitempty"`

		UserAgent string `json:"user_agent,omitempty"`
		APIKey    string `json:"api_key,omitempty"` // aka ClientID

		Status State `json:"status,omitempty"`

		Options map[string]string `json:"options,omitempty"` // holds all other values ...
	}
)

func (ds *Settings) Clone() Settings {
	s := Settings{
		Endpoint:  ds.Endpoint,
		UserAgent: ds.UserAgent,
		APIKey:    ds.APIKey,
		Status:    ds.Status,
	}

	if len(ds.Scopes) > 0 {
		s.Scopes = make([]string, len(ds.Scopes))
		copy(s.Scopes, ds.Scopes)
	}
	if len(ds.DefaultScopes) > 0 {
		s.DefaultScopes = make([]string, len(ds.DefaultScopes))
		copy(s.DefaultScopes, ds.DefaultScopes)
	}

	if ds.Credentials != nil {
		s.Credentials = ds.Credentials.Clone()
	}
	if len(ds.Options) > 0 {
		s.Options = make(map[string]string)
		for k, v := range ds.Options {
			s.Options[k] = v
		}
	}
	return s
}

// GetScopes returns the user-provided scopes, if set, or else falls back to the default scopes.
func (ds *Settings) GetScopes() []string {
	if len(ds.Scopes) > 0 {
		return ds.Scopes
	}
	return ds.DefaultScopes
}

// HasOption returns true if ds has a custom option opt.
func (ds *Settings) HasOption(opt string) bool {
	_, ok := ds.Options[opt]
	return ok
}

// GetOption returns the custom option opt if it exists or an empty string otherwise
func (ds *Settings) GetOption(opt string) string {
	if o, ok := ds.Options[opt]; ok {
		return o
	}
	return ""
}

// SetOptions registers a custom option o with key opt.
func (ds *Settings) SetOption(opt, o string) {
	if ds.Options == nil {
		ds.Options = make(map[string]string)
	}
	ds.Options[opt] = o
}
