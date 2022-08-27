package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/internal"
	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
)

type (
	localConfig struct {
		rootDir  string // the current working dir
		confDir  string // the fully qualified path to the conf dir
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
	}

	return c
}

func (c *localConfig) Name() string {
	return "simplecli"
}

func (c *localConfig) ShortName() string {
	return "sc"
}

func (c *localConfig) Copyright() string {
	return "Copyright 2022, transformative.services, https://txs.vc"
}

func (c *localConfig) About() string {
	return "a simple cli (sc) example"
}

func (c *localConfig) MajorVersion() int {
	return majorVersion
}

func (c *localConfig) MinorVersion() int {
	return minorVersion
}

func (c *localConfig) FixVersion() int {
	return fixVersion
}

//
//
//

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
	cs, err := internal.ReadSettingsFromFile(pathToFile) // FIXME: internal. will become an issue later
	if err != nil {
		cs = GetDefaultSettings()
		// save to the default location
		if err = cs.WriteToFile(pathToFile); err != nil {
			log.Fatal(err)
		}
	}

	// patch values from ENV if available

	// make it available for further call
	c.settings = cs
	return cs
}

func (c *localConfig) GetDefaultSettings() *settings.Settings {
	return &settings.Settings{
		Endpoint:      "http://localhost:8080",
		DefaultScopes: c.DefaultScopes(),
		Credentials:   &settings.Credentials{}, // add this to avoid NPEs further down
	}
}
