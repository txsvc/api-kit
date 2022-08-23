package config

import (
	"fmt"
	"path/filepath"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/internal"
	"github.com/txsvc/apikit/internal/settings"
)

const (
	// Version specifies the verion of the API and its structs
	apiVersion = "v1"
	// MajorVersion of the API
	majorVersion = 0
	// MinorVersion of the API
	minorVersion = 1
	// FixVersion of the API
	fixVersion = 2
)

type (
	GenericConfiguration struct {
		confDir string
	}
)

var (
	// interface guard to ensure that all required functions are implemented
	_ Configurator = (*GenericConfiguration)(nil)
)

func genericProvider() interface{} {
	return &GenericConfiguration{}
}

func VersionString() string {
	return fmt.Sprintf("%d.%d.%d", MajorVersion(), MinorVersion(), FixVersion())
}

func UserAgentString() string {
	return fmt.Sprintf("%s %d.%d.%d", ShortName(), MajorVersion(), MinorVersion(), FixVersion())
}

func ServerString() string {
	return fmt.Sprintf("%s %d.%d.%d", ShortName(), MajorVersion(), MinorVersion(), FixVersion())
}

func (c *GenericConfiguration) Name() string {
	return "appkit"
}

func (c *GenericConfiguration) ShortName() string {
	return "ak"
}

func (c *GenericConfiguration) Copyright() string {
	return "copyright 2022"
}

func (c *GenericConfiguration) About() string {
	return "about appkit"
}

func (c *GenericConfiguration) MajorVersion() int {
	return majorVersion
}

func (c *GenericConfiguration) MinorVersion() int {
	return minorVersion
}

func (c *GenericConfiguration) FixVersion() int {
	return fixVersion
}

func (c *GenericConfiguration) ApiVersion() string {
	return apiVersion
}

func (c *GenericConfiguration) DefaultConfigLocation() string {
	return DefaultConfigDirLocation
}

// GetConfigLocation returns the config location that was set using SetConfigLocation().
// If no location is defined, GetConfigLocation looks for ENV['CONFIG_LOCATION'] or
// returns DefaultConfigLocation() if no environment variable was set.
func (c *GenericConfiguration) GetConfigLocation() string {
	if len(c.confDir) == 0 {
		return stdlib.GetString(ConfigDirLocationENV, c.DefaultConfigLocation())
	}
	return c.confDir
}

func (c *GenericConfiguration) SetConfigLocation(loc string) {
	c.confDir = loc
}

func (c *GenericConfiguration) GetDefaultSettings() *settings.Settings {
	return &settings.Settings{
		Endpoint: "http://localhost:8080",
		DefaultScopes: []string{
			internal.ScopeApiRead,
			internal.ScopeApiWrite,
		},
	}
}

func (c *GenericConfiguration) GetSettings() *settings.Settings {

	pathToFile := filepath.Join(ResolveConfigLocation(), DefaultConfigFileName)

	s, err := internal.ReadSettingsFromFile(pathToFile)
	if err != nil {
		return c.GetDefaultSettings()
	}

	return s
}
