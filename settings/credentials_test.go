package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCredentialFile = "credentials.json"

func TestWriteReadCredentials(t *testing.T) {
	cred1 := &Credentials{
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
