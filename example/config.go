package example

import (
	"log"
	"os"
	"path/filepath"

	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/internal"
	"github.com/txsvc/apikit/internal/settings"
	"github.com/txsvc/stdlib/v2"
)

const (
	// Version specifies the verion of the API and its structs
	apiVersion = "v1"

	// MajorVersion of the API
	majorVersion = 0
	// MinorVersion of the API
	minorVersion = 1
	// FixVersion of the API
	fixVersion = 0
)

type (
	CmdLineConfiguration struct {
		rootDir string // the current working dir
		confDir string // the fully qualified path to the conf dir
	}
)

var (
	// interface guard to ensure that all required functions are implemented
	_ config.Configurator = (*CmdLineConfiguration)(nil)

	// curent client config
	currentSettings *settings.Settings
)

func NewExampleConfigProvider() interface{} {

	// get the current working dir. abort on error
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	c := &CmdLineConfiguration{
		rootDir: dir,
		confDir: "",
	}

	return c
}

func (c *CmdLineConfiguration) Name() string {
	return "simplecli"
}

func (c *CmdLineConfiguration) ShortName() string {
	return "sc"
}

func (c *CmdLineConfiguration) Copyright() string {
	return "Copyright 2022, transformative.services, https://txs.vc"
}

func (c *CmdLineConfiguration) About() string {
	return "a simple cli (sc) example"
}

func (c *CmdLineConfiguration) MajorVersion() int {
	return majorVersion
}

func (c *CmdLineConfiguration) MinorVersion() int {
	return minorVersion
}

func (c *CmdLineConfiguration) FixVersion() int {
	return fixVersion
}

func (c *CmdLineConfiguration) ApiVersion() string {
	return apiVersion
}

func (c *CmdLineConfiguration) DefaultConfigLocation() string {
	return config.DefaultConfigDirLocation
}

// GetConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *CmdLineConfiguration) GetConfigLocation() string {
	if len(c.confDir) == 0 {
		return stdlib.GetString(config.ConfigDirLocationENV, c.DefaultConfigLocation())
	}
	return c.confDir
}

func (c *CmdLineConfiguration) SetConfigLocation(loc string) {
	c.confDir = loc
}

func (c *CmdLineConfiguration) GetDefaultSettings() *settings.Settings {
	return &settings.Settings{
		Endpoint: "http://localhost:8080",
		DefaultScopes: []string{
			internal.ScopeApiRead,
			internal.ScopeApiWrite,
		},
	}
}

func (c *CmdLineConfiguration) GetSettings() *settings.Settings {
	if currentSettings != nil {
		return currentSettings
	}

	// try to load the dial settings
	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	cs, err := internal.ReadSettingsFromFile(pathToFile) // FIXME: internal. will become an issue later
	if err != nil {
		cs = config.GetDefaultSettings()
		// save to the default location
		if err = cs.WriteToFile(pathToFile); err != nil {
			log.Fatal(err)
		}
	}

	// patch values from ENV if available

	// make it available for further call
	currentSettings = cs
	return cs
}
