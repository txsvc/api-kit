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

	DefaultConfigName     = "config"
	DefaultConfigLocation = "./.config"
	DefaultEndpoint       = "http://localhost:8080" // only really useful for testing ...
)

type (
	// Info holds static information about a service or API
	Info struct {
		// name: the service's name in human-usable form
		name string
		// shortName: the abreviated version of the service's name
		shortName string
		// copyright: info on the copyright/owner of the service/api
		copyright string
		// about: a short description of the service/api
		about string
		// majorVersion: the major version of the service/api
		majorVersion int
		// minorVersion: the minor version of the service/api
		minorVersion int
		// fixVersion: the fix/patch version of the service/api
		fixVersion int
	}

	Configurator interface {
		// AppInfo returns static information about the app or service
		Info() *Info
		// GetScopes returns the user-provided scopes, if set, or else falls back to the default scopes.
		GetScopes() []string
		// GetConfigLocation returns the path to the config location, if set, or the default location otherwise.
		GetConfigLocation() string // './.config' unless explicitly set.
		// SetConfigLocation explicitly sets the location where the configuration is expected. The location's existence is NOT verified.
		SetConfigLocation(string)
		// Settings returns the app settings, if configured, or falls back to a default, minimal configuration
		Settings() *settings.DialSettings
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
	config_ interface{}
)

func init() {
	// makes sure that SOMETHING is initialized
	SetProvider(NewLocalConfigProvider())
}

func SetProvider(provider interface{}) {
	config_ = provider
}

func GetConfig() Configurator {
	return config_.(Configurator)
}

// SetConfigLocation sets the actual location without checking if the location actually exists !
func SetConfigLocation(loc string) {
	if config_ == nil {
		log.Fatal(ErrMissingConfigurator)
	}
	config_.(Configurator).SetConfigLocation(loc)
}

// ResolveConfigLocation returns the full path to the config location
func ResolveConfigLocation() string {
	cl := GetConfig().GetConfigLocation()
	if strings.HasPrefix(cl, ".") {
		// relative to working dir
		wd, err := os.Getwd()
		if err != nil {
			return GetConfig().GetConfigLocation()
		}
		return filepath.Join(wd, cl)
	}
	return cl
}
