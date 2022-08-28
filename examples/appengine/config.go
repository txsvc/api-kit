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

func (c *appConfig) Info() *config.Info {
	return c.info
}

func (c *appConfig) GetScopes() []string {
	if c.cfg_ != nil {
		return c.cfg_.GetScopes()
	}
	return defaultScopes()
}

// GetConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *appConfig) GetConfigLocation() string {
	if len(c.root) == 0 {
		return stdlib.GetString(config.ConfigDirLocationENV, config.DefaultConfigLocation)
	}
	return c.root
}

func (c *appConfig) SetConfigLocation(loc string) {
	c.root = loc
	if c.cfg_ != nil {
		c.cfg_ = nil // force a reload the next time GetSettings() is called ...
	}
}

func (c *appConfig) Settings() *settings.DialSettings {
	if c.cfg_ != nil {
		return c.cfg_
	}
	// make it available for future calls
	c.cfg_ = c.defaultSettings()
	return c.cfg_
}

func (c *appConfig) defaultSettings() *settings.DialSettings {
	cfg := settings.DialSettings{
		Endpoint:      config.DefaultEndpoint,
		DefaultScopes: defaultScopes(),
		Credentials:   &settings.Credentials{}, // add this to avoid NPEs further down
	}
	// patch values from ENV, if available
	cfg.Endpoint = stdlib.GetString(config.APIEndpointENV, cfg.Endpoint)
	return &cfg
}

func defaultScopes() []string {
	// FIXME: this gives basic read access to the API. Is this what we want?
	return []string{
		auth.ScopeApiRead,
		auth.ScopeApiWrite,
	}
}
