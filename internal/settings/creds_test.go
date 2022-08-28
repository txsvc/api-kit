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
