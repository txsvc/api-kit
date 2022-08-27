package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRange(t *testing.T) {

	offset, length := ParseRange("")
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), length)

	offset, length = ParseRange("0-100") // not valid
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), length)

	offset, length = ParseRange("bytes=10-100") // valid, single range only
	assert.Equal(t, int64(10), offset)
	assert.Equal(t, int64(90), length)

	offset, length = ParseRange("bytes=0-50, 100-150") // not valid, no support for multipart ranges
	assert.Equal(t, int64(0), offset)
	assert.Equal(t, int64(-1), length)
}
