package helpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	s := RandString(10)
	assert.NotEmpty(t, s)
	assert.Equal(t, 10, len(s))

	s = RandString(1024)
	assert.NotEmpty(t, s)
	assert.Equal(t, 1024, len(s))

	// just make this long enough that it is not likely that anything is missing. could randomly fail though.

	containsAllChars := true
	for i := 0; i < len(letterBytes); i++ {
		if !strings.Contains(s, string(letterBytes[i])) {
			containsAllChars = false
			fmt.Printf("%s: missing: '%s'\n", s, string(letterBytes[i]))
			break
		}
	}
	assert.True(t, containsAllChars)
}

func TestRandStringSimple(t *testing.T) {
	s := RandStringSimple(16)
	assert.NotEmpty(t, s)
	assert.Equal(t, 16, len(s))

	s = RandStringSimple(64)
	assert.NotEmpty(t, s)
	assert.Equal(t, 64, len(s))

	fmt.Println(s)
}

func TestRandPasswordString(t *testing.T) {
	s := RandPasswordString(16)
	assert.NotEmpty(t, s)
	assert.Equal(t, 16, len(s))

	s = RandPasswordString(64)
	assert.NotEmpty(t, s)
	assert.Equal(t, 64, len(s))

	fmt.Println(s)
}
