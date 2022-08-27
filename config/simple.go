package config

import (
	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
)

type (
	simpleConfig struct {
		root string // the fully qualified path to the conf dir
	}
)

var (
	// interface guard to ensure that all required functions are implemented
	_ Configurator = (*simpleConfig)(nil)
)

func NewSimpleConfigProvider() interface{} {
	return &simpleConfig{}
}

func (c *simpleConfig) Name() string {
	return "appkit"
}

func (c *simpleConfig) ShortName() string {
	return "appkit"
}

func (c *simpleConfig) Copyright() string {
	return "copyright 2022, transformative.services, https://txs.vc"
}

func (c *simpleConfig) About() string {
	return "about appkit"
}

func (c *simpleConfig) MajorVersion() int {
	return majorVersion
}

func (c *simpleConfig) MinorVersion() int {
	return minorVersion
}

func (c *simpleConfig) FixVersion() int {
	return fixVersion
}

//
//
//

func (c *simpleConfig) DefaultScopes() []string {
	return []string{
		auth.ScopeApiRead,
		auth.ScopeApiWrite,
	}
}

// GetConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *simpleConfig) GetConfigLocation() string {
	if len(c.root) == 0 {
		return stdlib.GetString(ConfigDirLocationENV, c.DefaultConfigLocation())
	}
	return c.root
}

func (c *simpleConfig) SetConfigLocation(loc string) {
	c.root = loc
}

func (c *simpleConfig) DefaultConfigLocation() string {
	return DefaultConfigDirLocation
}

func (c *simpleConfig) GetSettings() *settings.Settings {
	return c.GetDefaultSettings()
}

func (c *simpleConfig) GetDefaultSettings() *settings.Settings {
	return &settings.Settings{
		Endpoint:      "http://localhost:8080",
		DefaultScopes: c.DefaultScopes(),
		Credentials:   &settings.Credentials{}, // add this to avoid NPEs further down
	}
}
