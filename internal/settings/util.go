package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ReadSettingsFromFile(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ds := Settings{}
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func ReadCredentialsFromFile(path string) (*Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cred := Credentials{}
	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

func WriteSettingsToFile(cfg *Settings, path string) error {
	buf, err := json.MarshalIndent(cfg, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, buf, filePerm)
}

func WriteCredentialsToFile(cred *Credentials, path string) error {
	buf, err := json.MarshalIndent(cred, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	return os.WriteFile(path, buf, filePerm)
}
