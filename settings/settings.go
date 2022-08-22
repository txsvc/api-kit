// Most of the code is lifted from
// https://github.com/googleapis/google-api-go-client/blob/main/internal/settings.go
//
// For details and copyright etc. see above url.
//
package settings

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	indentChar             = " "
	filePerm   fs.FileMode = 0644
)

type (
	// DialSettings holds information needed to establish a connection with a
	// backend API service or to simply configure some code.
	DialSettings struct {
		Endpoint        string `json:"endpoint,omitempty"`
		DefaultEndpoint string `json:"default_endpoint,omitempty"`

		Scopes        []string `json:"scopes,omitempty"`
		DefaultScopes []string `json:"default_scopes,omitempty"`

		Credentials         *Credentials `json:"credentials,omitempty"`
		InternalCredentials *Credentials `json:"internal_credentials,omitempty"`
		CredentialsFile     string       `json:"credentials_file,omitempty"`

		Options map[string]string `json:"options,omitempty"` // holds all other options ...

		UserAgent string `json:"user_agent,omitempty"`
		APIKey    string `json:"api_key,omitempty"` // aka ClientID

		NoAuth         bool `json:"no_auth,omitempty"`
		SkipValidation bool `json:"skip_validation,omitempty"`
	}
)

// GetScopes returns the user-provided scopes, if set, or else falls back to the default scopes.
func (ds *DialSettings) GetScopes() []string {
	if len(ds.Scopes) > 0 {
		return ds.Scopes
	}
	return ds.DefaultScopes
}

// HasOption returns true if ds has a custom option opt.
func (ds *DialSettings) HasOption(opt string) bool {
	_, ok := ds.Options[opt]
	return ok
}

// GetOption returns the custom option opt if it exists or an empty string otherwise
func (ds *DialSettings) GetOption(opt string) string {
	if o, ok := ds.Options[opt]; ok {
		return o
	}
	return ""
}

// SetOptions registers a custom option o with key opt.
func (ds *DialSettings) SetOption(opt, o string) {
	if ds.Options == nil {
		ds.Options = make(map[string]string)
	}
	ds.Options[opt] = o
}

// Validate reports an error if ds is invalid.
func (ds *DialSettings) Validate() error {
	if ds.SkipValidation {
		return nil
	}

	hasCreds := ds.APIKey != "" || ds.CredentialsFile != "" || ds.Credentials != nil
	if ds.NoAuth && hasCreds {
		return errors.New("settings.WithoutAuthentication is incompatible with any option that provides credentials")
	}

	nCreds := 0
	if ds.APIKey != "" {
		nCreds++
	}
	if ds.CredentialsFile != "" {
		nCreds++
	}
	if ds.Credentials != nil {
		nCreds++
	}
	if nCreds > 1 {
		return errors.New("multiple credential options provided")
	}

	return nil
}

func (ds *DialSettings) WriteToFile(path string) error {
	cfg, err := json.MarshalIndent(ds, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return ioutil.WriteFile(path, cfg, filePerm)
}

func ReadSettingsFromFile(path string) (*DialSettings, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ds := DialSettings{}
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}
