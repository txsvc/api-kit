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

func ReadSettingsFromFile(path string) (*settings.Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ds := settings.Settings{}
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func WriteSettingsToFile(cfg *settings.Settings, path string) error {
	buf, err := json.MarshalIndent(cfg, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, buf, filePerm)
}

func ReadCredentialsFromFile(path string) (*settings.Credentials, error) {
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

func WriteCredentialsToFile(cred *settings.Credentials, path string) error {
	buf, err := json.MarshalIndent(cred, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, buf, filePerm)
}
