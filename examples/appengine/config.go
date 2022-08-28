package main

import (
	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
)

// FIXME: make this Google AppEngine specific !

var (
	// interface guard to ensure that all required functions are implemented
	_ config.Configurator = (*appConfig)(nil)
)

func (c *appConfig) AppInfo() *config.Info {
	return c.info
}

func (c *appConfig) DefaultScopes() []string {
	return []string{
		auth.ScopeApiRead,
		auth.ScopeApiWrite,
	}
}

// GetConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *appConfig) GetConfigLocation() string {
	if len(c.root) == 0 {
		return stdlib.GetString(config.ConfigDirLocationENV, c.DefaultConfigLocation())
	}
	return c.root
}

func (c *appConfig) SetConfigLocation(loc string) {
	c.root = loc
}

func (c *appConfig) DefaultConfigLocation() string {
	return config.DefaultConfigDirLocation
}

func (c *appConfig) GetSettings() *settings.Settings {
	return c.GetDefaultSettings()
}

func (c *appConfig) GetDefaultSettings() *settings.Settings {
	return &settings.Settings{
		Endpoint:      "http://localhost:8080",
		DefaultScopes: c.DefaultScopes(),
		Credentials:   &settings.Credentials{}, // add this to avoid NPEs further down
	}
}
