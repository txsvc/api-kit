package apikit

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	// ShutdownDelay is the time to wait for all request, go-routines etc complete
	ShutdownDelay = 10 // seconds
)

type (

	// SetupFunc creates a new, fully configured mux
	SetupFunc func() *echo.Echo
	// ShutdownFunc is called before the server stops
	ShutdownFunc func(*App)

	// app holds all configs for the listener
	App struct {
		mux              *echo.Echo
		shutdown         ShutdownFunc
		errorHandlerImpl echo.HTTPErrorHandler
		// other settings
		logLevel      log.Lvl
		shutdownDelay time.Duration
	}
)

func New(setupFunc SetupFunc, shutdownFunc ShutdownFunc, errorHandler echo.HTTPErrorHandler) (*App, error) {
	if setupFunc == nil || shutdownFunc == nil {
		return nil, ErrInvalidConfiguration
	}

	app := &App{
		shutdown:         shutdownFunc,
		errorHandlerImpl: errorHandler,
		logLevel:         log.INFO,
		shutdownDelay:    ShutdownDelay * time.Second,
	}

	app.mux = setupFunc()
	if app.mux == nil {
		return nil, ErrInvalidConfiguration
	}

	return app, nil
}
