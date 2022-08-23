package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopes(t *testing.T) {
	cfg1 := Settings{
		DefaultScopes: []string{"a", "b"},
	}
	assert.NotEmpty(t, cfg1.GetScopes())

	cfg2 := Settings{
		Scopes: []string{"A", "B"},
	}
	assert.NotEmpty(t, cfg2.GetScopes())
}

func TestOptions(t *testing.T) {
	cfg1 := Settings{}
	assert.Nil(t, cfg1.Options)
	assert.False(t, cfg1.HasOption("FOO"))

	opt := cfg1.GetOption("FOO")
	assert.Empty(t, opt)

	cfg1.SetOption("FOO", "x")
	assert.True(t, cfg1.HasOption("FOO"))
	opt = cfg1.GetOption("FOO")
	assert.Equal(t, "x", opt)
}
