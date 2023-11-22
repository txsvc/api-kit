package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/stdlib/v2/settings"
)

const testCredentialFile = "test.json"

func TestWriteReadSettings(t *testing.T) {
	settings1 := &settings.DialSettings{
		Endpoint: "x",
		//DefaultEndpoint: "X",
		Scopes:        []string{"a", "b"},
		DefaultScopes: []string{"A", "B"},
		UserAgent:     "agent",
	}
	settings1.SetOption("FOO", "x")
	settings1.SetOption("BAR", "x")

	err := WriteDialSettings(settings1, testCredentialFile)
	assert.NoError(t, err)

	settings2, err := ReadDialSettings(testCredentialFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, settings2)
	assert.Equal(t, settings1, settings2)

	// cleanup
	os.Remove(testCredentialFile)
}

func TestWriteReadCredentials(t *testing.T) {
	cred1 := &settings.Credentials{
		ProjectID: "project",
		ClientID:  "client",
		Token:     "token",
		Expires:   42,
	}

	err := WriteCredentials(cred1, testCredentialFile)
	assert.NoError(t, err)

	cred2, err := ReadCredentials(testCredentialFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, cred2)
	assert.Equal(t, cred1, cred2)

	// cleanup
	os.Remove(testCredentialFile)
}
