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
// i.e. there should be a version 'v0.1.0' !
const (
	majorVersion = 0
	minorVersion = 1
	fixVersion   = 0
)

type (
	localConfig struct {
		// app info
		info *Info
		// path to configuration settings
		rootDir string // the current working dir
		confDir string // the fully qualified path to the conf dir
		// cached settings
		settings *settings.Settings
	}
)

var (
	// interface guard to ensure that all required functions are implemented
	_ Configurator = (*localConfig)(nil)
)

func NewLocalConfigProvider() interface{} {

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

func (c *localConfig) AppInfo() *Info {
	return c.info
}

func (c *localConfig) DefaultScopes() []string {
	return []string{
		auth.ScopeApiRead,
		auth.ScopeApiWrite,
	}
}

// GetConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *localConfig) GetConfigLocation() string {
	if len(c.confDir) == 0 {
		return stdlib.GetString(ConfigDirLocationENV, c.DefaultConfigLocation())
	}
	return c.confDir
}

func (c *localConfig) SetConfigLocation(loc string) {
	c.confDir = loc
}

func (c *localConfig) DefaultConfigLocation() string {
	return DefaultConfigDirLocation
}

func (c *localConfig) GetSettings() *settings.Settings {
	if c.settings != nil {
		return c.settings
	}

	// try to load the dial settings
	pathToFile := filepath.Join(ResolveConfigLocation(), DefaultConfigFileName)
	cfg, err := helpers.ReadSettingsFromFile(pathToFile)
	if err != nil {
		cfg = GetDefaultSettings()
		// save to the default location
		if err = helpers.WriteSettingsToFile(cfg, pathToFile); err != nil {
			log.Fatal(err)
		}
	}

	// patch values from ENV if available

	// make it available for further call
	c.settings = cfg
	return cfg
}

func (c *localConfig) GetDefaultSettings() *settings.Settings {
	return &settings.Settings{
		Endpoint:      "http://localhost:8080",
		DefaultScopes: c.DefaultScopes(),
		Credentials:   &settings.Credentials{}, // add this to avoid NPEs further down
	}
}
