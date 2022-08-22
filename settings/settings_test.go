package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	valid := []DialSettings{
		{},
		{APIKey: "x"},
		{CredentialsFile: "f"},
		{Credentials: &Credentials{}},
	}

	for _, ds := range valid {
		err := ds.Validate()
		assert.NoError(t, err)
	}

	notValid := []DialSettings{
		{NoAuth: true, APIKey: "x"},
		{NoAuth: true, CredentialsFile: "f"},
		{NoAuth: true, Credentials: &Credentials{}},
		{APIKey: "x", CredentialsFile: "f"},
		{CredentialsFile: "f", Credentials: &Credentials{}},
		{APIKey: "x", CredentialsFile: "f", Credentials: &Credentials{}},
	}

	for _, ds := range notValid {
		err := ds.Validate()
		assert.Error(t, err)
	}
}

func TestScopes(t *testing.T) {
	cfg1 := DialSettings{
		DefaultScopes: []string{"a", "b"},
	}
	assert.NotEmpty(t, cfg1.GetScopes())

	cfg2 := DialSettings{
		Scopes: []string{"A", "B"},
	}
	assert.NotEmpty(t, cfg2.GetScopes())
}

func TestOptions(t *testing.T) {
	cfg1 := DialSettings{}
	assert.Nil(t, cfg1.Options)
	assert.False(t, cfg1.HasOption("FOO"))

	opt := cfg1.GetOption("FOO")
	assert.Empty(t, opt)

	cfg1.SetOption("FOO", "x")
	assert.True(t, cfg1.HasOption("FOO"))
	opt = cfg1.GetOption("FOO")
	assert.Equal(t, "x", opt)
}

func TestWriteReadSettings(t *testing.T) {
	settings1 := &DialSettings{
		Endpoint:        "x",
		DefaultEndpoint: "X",
		Scopes:          []string{"a", "b"},
		DefaultScopes:   []string{"A", "B"},
		UserAgent:       "agent",
		APIKey:          "api_key",
	}
	settings1.SetOption("FOO", "x")
	settings1.SetOption("BAR", "x")

	err := settings1.WriteToFile(testCredentialFile)
	assert.NoError(t, err)

	settings2, err := ReadSettingsFromFile(testCredentialFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, settings2)
	assert.Equal(t, settings1, settings2)

	// cleanup
	os.Remove(testCredentialFile)
}
