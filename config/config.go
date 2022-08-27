package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/txsvc/apikit/internal/settings"
)

const (
	DefaultConfigDirLocation = "./.config"
	DefaultConfigFileName    = "config"

	ConfigDirLocationENV = "CONFIG_LOCATION"
	APIEndpointENV       = "API_ENDPOINT"
)

type (
	ConfigProviderFunc func() interface{}

	Configurator interface {
		Name() string      // name of the project / the real etc
		ShortName() string // abreviated name, used for e.g. the cli tool
		Copyright() string
		About() string

		MajorVersion() int
		MinorVersion() int
		FixVersion() int

		DefaultConfigLocation() string // default: ./.config

		GetConfigLocation() string // same as DefaultConfigLocation() unless explicitly set
		SetConfigLocation(string)

		// client & endpoint settings and credentials
		GetDefaultSettings() *settings.Settings
		GetSettings() *settings.Settings
	}
)

var (
	// ErrMissingConfigurator indicates that the config package is not initialized
	ErrMissingConfigurator = errors.New("missing configurator")
	// ErrInitializingConfiguration indicates that the client could not be initialized
	ErrInitializingConfiguration = errors.New("error initializing configuration")
	// ErrInvalidConfiguration indicates that parameters used to configure the service were invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")

	confProvider interface{}
)

func init() {
	// makes sure that SOMETHING is initialized
	InitConfigProvider(NewSimpleConfigProvider())
}

func InitConfigProvider(provider interface{}) {
	confProvider = provider
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

func Name() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).Name()
}

func ShortName() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).ShortName()
}

func Copyright() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).Copyright()
}

func About() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).About()
}

func MajorVersion() int {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).MajorVersion()
}

func MinorVersion() int {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).MinorVersion()
}
func FixVersion() int {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).FixVersion()
}

// DefaultConfigLocation returns a default location e.g. %HOME/.config
func DefaultConfigLocation() string {
	if confProvider == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	return confProvider.(Configurator).DefaultConfigLocation()
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
