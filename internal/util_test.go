package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/apikit/internal/settings"
)

const testCredentialFile = "test.json"

func TestWriteReadSettings(t *testing.T) {
	settings1 := &settings.Settings{
		Endpoint: "x",
		//DefaultEndpoint: "X",
		Scopes:        []string{"a", "b"},
		DefaultScopes: []string{"A", "B"},
		UserAgent:     "agent",
		APIKey:        "api_key",
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

func TestWriteReadCredentials(t *testing.T) {
	cred1 := &settings.Credentials{
		ProjectID: "project",
		UserID:    "user",
		Token:     "token",
		Expires:   42,
	}

	err := cred1.WriteToFile(testCredentialFile)
	assert.NoError(t, err)

	cred2, err := ReadCredentialsFromFile(testCredentialFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, cred2)
	assert.Equal(t, cred1, cred2)

	// cleanup
	os.Remove(testCredentialFile)
}
