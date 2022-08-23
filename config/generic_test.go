package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionStrings(t *testing.T) {
	assert.NotEmpty(t, VersionString())
	assert.NotEmpty(t, UserAgentString())
	assert.NotEmpty(t, ServerString())
}

func TestGenericConfig(t *testing.T) {
	conf := genericProvider().(*GenericConfiguration)
	assert.NotNil(t, conf)

	assert.Equal(t, conf.Name(), Name())
	assert.Equal(t, conf.ShortName(), ShortName())
	assert.Equal(t, conf.Copyright(), Copyright())
	assert.Equal(t, conf.About(), About())

	assert.Equal(t, conf.MajorVersion(), MajorVersion())
	assert.Equal(t, conf.MinorVersion(), MinorVersion())
	assert.Equal(t, conf.FixVersion(), FixVersion())
	assert.Equal(t, conf.ApiVersion(), ApiVersion())

	assert.Equal(t, conf.DefaultConfigLocation(), DefaultConfigLocation())
	assert.Equal(t, conf.GetConfigLocation(), GetConfigLocation())
	assert.Equal(t, conf.GetConfigLocation(), conf.DefaultConfigLocation())
}

func TestSetConfigLocation(t *testing.T) {
	conf := genericProvider().(*GenericConfiguration)
	assert.NotNil(t, conf)

	InitConfigProvider(conf)

	assert.Equal(t, DefaultConfigDirLocation, DefaultConfigLocation())
	assert.Equal(t, DefaultConfigDirLocation, GetConfigLocation())
	assert.Equal(t, conf.DefaultConfigLocation(), conf.GetConfigLocation())

	conf.SetConfigLocation("$HOME/.config")

	assert.Equal(t, "$HOME/.config", GetConfigLocation())
	assert.Equal(t, DefaultConfigDirLocation, DefaultConfigLocation())
}

func TestGetDefaultSettings(t *testing.T) {
	conf := genericProvider().(*GenericConfiguration)
	assert.NotNil(t, conf)

	InitConfigProvider(conf)
	ds := GetDefaultSettings()

	assert.NotNil(t, ds)
	assert.NotEmpty(t, ds)
}

func TestGetSettings(t *testing.T) {
	conf := genericProvider().(*GenericConfiguration)
	assert.NotNil(t, conf)

	InitConfigProvider(conf)
	ds := GetSettings()

	assert.NotNil(t, ds)
	assert.NotEmpty(t, ds)
}
