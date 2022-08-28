package helpers

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/txsvc/apikit/internal/settings"
)

const (
	indentChar             = "  "
	filePerm   fs.FileMode = 0644
)

func ReadDialSettings(path string) (*settings.DialSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ds := settings.DialSettings{}
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func WriteDialSettings(cfg *settings.DialSettings, path string) error {
	buf, err := json.MarshalIndent(cfg, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, buf, filePerm)
}

func ReadCredentials(path string) (*settings.Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cred := settings.Credentials{}
	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

func WriteCredentials(cred *settings.Credentials, path string) error {
	buf, err := json.MarshalIndent(cred, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, buf, filePerm)
}
