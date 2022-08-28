package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneCredentials(t *testing.T) {
	cred := Credentials{
		ProjectID: "p",
		UserID:    "u",
		Token:     "t",
		Expires:   10,
	}
	dup := cred.Clone()
	assert.Equal(t, &cred, dup)
}

func TestValidation(t *testing.T) {
	cred1 := Credentials{}
	assert.False(t, cred1.IsValid())

	cred2 := Credentials{
		ProjectID: "p",
		UserID:    "u",
		Token:     "t",
		Expires:   10, // forces a fail ...
	}
	assert.False(t, cred2.IsValid())

	cred2.Expires = 0
	assert.True(t, cred2.IsValid())
}

func TestExpiration(t *testing.T) {
	cred := Credentials{
		ProjectID: "p",
		UserID:    "u",
		Token:     "t",
		Expires:   10,
	}
	assert.True(t, cred.Expired())

	cred.Expires = 0
	assert.False(t, cred.Expired())
}
