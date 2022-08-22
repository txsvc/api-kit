// Most of the code is lifted from
// https://cs.opensource.google/go/x/oauth2/+/d3ed0bb2:google/default.go
//
// For details and copyright etc. see above url.
//
package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Credentials struct {
	ProjectID string `json:"project_id,omitempty"` // may be empty
	UserID    string `json:"user_id,omitempty"`    // may be empty
	Token     string `json:"token,omitempty"`      // may be empty
	Expires   int64  `json:"expires,omitempty"`    // 0 = never
}

func (cred *Credentials) WriteToFile(path string) error {
	cfg, err := json.MarshalIndent(cred, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return ioutil.WriteFile(path, cfg, filePerm)
}

func ReadCredentialsFromFile(path string) (*Credentials, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cred := Credentials{}
	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}
