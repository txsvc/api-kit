package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionStrings(t *testing.T) {
	conf := NewLocalConfigProvider().(*localConfig)
	assert.NotNil(t, conf)

	info := conf.Info()
	assert.NotNil(t, info)
	assert.NotEmpty(t, info)

	assert.Equal(t, majorVersion, info.MajorVersion())
	assert.Equal(t, minorVersion, info.MinorVersion())
	assert.Equal(t, fixVersion, info.FixVersion())
	assert.NotEmpty(t, info.Name())
	assert.NotEmpty(t, info.ShortName())
	assert.NotEmpty(t, info.Copyright())
	assert.NotEmpty(t, info.About())
}
