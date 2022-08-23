// Most of the code is lifted from
// https://github.com/googleapis/google-api-go-client/blob/main/internal/settings.go
//
// For details and copyright etc. see above url.
package settings

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/txsvc/stdlib/v2"
)

const (
	indentChar             = "  "
	filePerm   fs.FileMode = 0644
)

type (
	// Settings holds information needed to establish a connection with a
	// backend API service or to simply configure some code.
	Settings struct {
		Endpoint string `json:"endpoint,omitempty"`

		Scopes        []string `json:"scopes,omitempty"`
		DefaultScopes []string `json:"default_scopes,omitempty"`

		Credentials *Credentials `json:"credentials,omitempty"`
		//InternalCredentials *Credentials `json:"internal_credentials,omitempty"`
		//CredentialsFile     string       `json:"credentials_file,omitempty"`
		//NoAuth              bool         `json:"no_auth,omitempty"`

		UserAgent string `json:"user_agent,omitempty"`
		APIKey    string `json:"api_key,omitempty"` // aka ClientID

		Options map[string]string `json:"options,omitempty"` // holds all other values ...
	}

	Credentials struct {
		ProjectID string `json:"project_id,omitempty"` // may be empty
		UserID    string `json:"user_id,omitempty"`    // may be empty, aka client_id
		Token     string `json:"token,omitempty"`      // may be empty
		Expires   int64  `json:"expires,omitempty"`    // 0 = never, < 0 = invalid, > 0 = unix timestamp
	}
)

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

func (ds *Settings) WriteToFile(path string) error {
	cfg, err := json.MarshalIndent(ds, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, cfg, filePerm)
}

// IsValid test if Crendentials is valid
func (c *Credentials) IsValid() bool {
	// explicitly set to invalid
	if c.Expires < 0 {
		return false
	}
	// attributes must be set
	if len(c.Token) == 0 || len(c.ProjectID) == 0 || len(c.UserID) == 0 {
		return false
	}
	// never expires
	if c.Expires == 0 {
		return true
	}
	// check if Expires is still in the future
	return c.Expires > stdlib.Now()
}

func (cred *Credentials) WriteToFile(path string) error {
	cfg, err := json.MarshalIndent(cred, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, cfg, filePerm)
}