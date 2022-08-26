package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
	assert.Equal(t, Info, l.level) // assumes ENV['LOG_LEVEL'] is not set

	l = NewLogger(os.Stdout, "")
	assert.NotNil(t, l)
	assert.Equal(t, Info, l.level) // assumes ENV['LOG_LEVEL'] is not set

	l = NewLogger(os.Stdout, "debug")
	assert.NotNil(t, l)
	assert.Equal(t, Debug, l.level)

	l = NewLogger(os.Stdout, "warn")
	assert.NotNil(t, l)
	assert.Equal(t, Warn, l.level)

	l = NewLogger(os.Stdout, "error")
	assert.NotNil(t, l)
	assert.Equal(t, Error, l.level)
}

func TestLoggerLog(t *testing.T) {
	l := NewLogger(os.Stdout, "")
	l.Info("this is INFO on INFO")
	l.Warn("this is WARN on INFO")
	l.Error("this is ERROR on INFO")
	l.Debug("this is DEBUG on INFO")
	fmt.Println("")

	l = NewLogger(os.Stdout, "warn")
	l.Info("this is INFO on WARN")
	l.Warn("this is WARN on WARN")
	l.Error("this is ERROR on WARN")
	l.Debug("this is DEBUG on WARN")
	fmt.Println("")

	l = NewLogger(os.Stdout, "debug")
	l.Info("this is INFO on DEBUG")
	l.Warn("this is WARN on DEBUG")
	l.Error("this is ERROR on DEBUG")
	l.Debug("this is DEBUG on DEBUG")
	fmt.Println("")

	l = NewLogger(os.Stdout, "error")
	l.Info("this is INFO on ERROR")
	l.Warn("this is WARN on ERROR")
	l.Error("this is ERROR on ERROR")
	l.Debug("this is DEBUG on ERROR")
}
