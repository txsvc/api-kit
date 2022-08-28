package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/txsvc/apikit/internal/settings"
)

const (
	// runtime settings
	PortEnv        = "PORT"
	APIEndpointENV = "API_ENDPOINT"
	// client settings
	ForceTraceEnv = "APIKIT_FORCE_TRACE"
	// config settings
	ConfigDirLocationENV = "CONFIG_LOCATION"

	DefaultConfigDirLocation = "./.config"
	DefaultConfigFileName    = "config"
)

type (
	Info struct {
		name         string
		shortName    string
		copyright    string
		about        string
		majorVersion int
		minorVersion int
		fixVersion   int
	}

	ConfigProviderFunc func() interface{}

	Configurator interface {
		AppInfo() *Info

		DefaultScopes() []string

		GetConfigLocation() string // same as DefaultConfigLocation() unless explicitly set
		SetConfigLocation(string)
		DefaultConfigLocation() string // default: ./.config

		// client & endpoint settings and credentials
		GetSettings() *settings.Settings
		GetDefaultSettings() *settings.Settings
	}
)

var (
	// ErrMissingConfigurator indicates that the config package is not initialized
	ErrMissingConfigurator = errors.New("missing configurator")
	// ErrInitializingConfiguration indicates that the client could not be initialized
	ErrInitializingConfiguration = errors.New("error initializing configuration")
	// ErrInvalidConfiguration indicates that parameters used to configure the service were invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// the config "singleton"
	confProvider interface{}
)

func init() {
	// makes sure that SOMETHING is initialized
	InitConfigProvider(NewLocalConfigProvider())
}

func InitConfigProvider(provider interface{}) {
	confProvider = provider
}

//
// "static" functions to match Configurator interface
//

func AppInfo() *Info {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).AppInfo()
}

func GetDefaultScopes() []string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).DefaultScopes()
}

// ConfigLocation returns the actual location or DefaultConfigLocation() if undefined
func GetConfigLocation() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).GetConfigLocation()
}

// SetConfigLocation sets the actual location without checking if the location actually exists !
func SetConfigLocation(loc string) {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	confProvider.(Configurator).SetConfigLocation(loc)
}

// DefaultConfigLocation returns a default location e.g. %HOME/.config
func DefaultConfigLocation() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).DefaultConfigLocation()
}

// ResolveConfigLocation returns the full path to the config location
func ResolveConfigLocation() string {
	cl := GetConfigLocation()
	if strings.HasPrefix(cl, ".") {
		// relative to working dir
		wd, err := os.Getwd()
		if err != nil {
			return DefaultConfigLocation()
		}
		return filepath.Join(wd, cl)
	}
	return GetConfigLocation()
}

func GetDefaultSettings() *settings.Settings {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).GetDefaultSettings()
}

func GetSettings() *settings.Settings {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).GetSettings()
}
