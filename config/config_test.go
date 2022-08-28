package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	conf1 := GetConfig()
	assert.NotNil(t, conf1)

	conf2 := NewLocalConfigProvider().(*localConfig)
	assert.NotNil(t, conf2)
}

func TestSetConfig(t *testing.T) {
	conf := NewLocalConfigProvider().(*localConfig)
	assert.NotNil(t, conf)

	SetProvider(conf)
	assert.Equal(t, conf, GetConfig())
}
