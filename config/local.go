package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/helpers"
	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
)

// the below version numbers should match the git release tags,
// i.e. there should be a version 'v0.1.0' on branch main !
const (
	majorVersion = 0
	minorVersion = 1
	fixVersion   = 0
)

type (
	localConfig struct {
		// the interface to implement
		Configurator

		// app info
		info *Info
		// path to configuration settings
		rootDir string // the current working dir
		confDir string // the fully qualified path to the conf dir
		// cached settings
		cfg_ *settings.DialSettings
	}
)

func NewLocalConfigProvider() Configurator {

	// get the current working dir. abort on error
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	c := &localConfig{
		rootDir: dir,
		confDir: "",
		info: &Info{
			name:         "appkit",
			shortName:    "appkit",
			copyright:    "Copyright 2022, transformative.services, https://txs.vc",
			about:        "about appkit",
			majorVersion: majorVersion,
			minorVersion: minorVersion,
			fixVersion:   fixVersion,
		},
	}

	return c
}

func (c *localConfig) Info() *Info {
	return c.info
}

func (c *localConfig) GetScopes() []string {
	if c.cfg_ != nil {
		return c.cfg_.GetScopes()
	}
	return defaultScopes()
}

// ConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *localConfig) ConfigLocation() string {
	if len(c.confDir) == 0 {
		return stdlib.GetString(ConfigDirLocationENV, DefaultConfigLocation)
	}
	return c.confDir
}

func (c *localConfig) SetConfigLocation(loc string) {
	c.confDir = loc
	if c.cfg_ != nil {
		c.cfg_ = nil // force a reload the next time GetSettings() is called ...
	}
}

func (c *localConfig) Settings() *settings.DialSettings {
	if c.cfg_ != nil {
		return c.cfg_
	}

	// try to load the dial settings
	pathToFile := filepath.Join(c.ConfigLocation(), DefaultConfigName)
	cfg, err := helpers.ReadDialSettings(pathToFile)
	if err != nil {
		cfg = c.defaultSettings()
		// save to the default location
		if err = helpers.WriteDialSettings(cfg, pathToFile); err != nil {
			log.Fatal(err)
		}
	}

	// patch values from ENV, if available
	cfg.Endpoint = stdlib.GetString(APIEndpointENV, cfg.Endpoint)

	// make it available for future calls
	c.cfg_ = cfg
	return c.cfg_
}

func (c *localConfig) defaultSettings() *settings.DialSettings {
	return &settings.DialSettings{
		Endpoint:      DefaultEndpoint,
		DefaultScopes: defaultScopes(),
		Credentials:   &settings.Credentials{}, // add this to avoid NPEs further down
	}
}

func defaultScopes() []string {
	// FIXME: this gives basic read access to the API. Is this what we want?
	return []string{
		auth.ScopeApiRead,
	}
}
