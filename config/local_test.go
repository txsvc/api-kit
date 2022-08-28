package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLocalProvider(t *testing.T) {
	SetProvider(NewLocalConfigProvider())

	cfg := GetConfig()
	assert.NotNil(t, cfg)

	assert.NotNil(t, cfg.Info())
	assert.NotNil(t, cfg.Settings())
	assert.NotEmpty(t, cfg.GetScopes())
}

func TestConfigLocation(t *testing.T) {
	SetProvider(NewLocalConfigProvider())

	cfg := GetConfig()
	assert.Equal(t, cfg, GetConfig())

	path := cfg.GetConfigLocation()
	assert.NotEmpty(t, path)
	assert.Equal(t, DefaultConfigLocation, path)

	cfg.SetConfigLocation("$HOME/.config")
	assert.Equal(t, "$HOME/.config", cfg.GetConfigLocation())
}

func TestGetSettings(t *testing.T) {
	conf := NewLocalConfigProvider().(*localConfig)
	assert.NotNil(t, conf)

	SetProvider(conf)
	ds := GetConfig().Settings()
	assert.NotNil(t, ds)
	assert.NotEmpty(t, ds)
}
