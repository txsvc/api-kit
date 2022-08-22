package apikit

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	// default ports to listen on
	PORT_DEFAULT     = "8080"
	PORT_DEFAULT_TLS = "433"

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
		root          string
	}
)

func New(setupFunc SetupFunc, shutdownFunc ShutdownFunc, errorHandler echo.HTTPErrorHandler) (*App, error) {
	if setupFunc == nil || shutdownFunc == nil {
		return nil, ErrInvalidConfiguration
	}

	app := &App{
		mux:              setupFunc(),
		shutdown:         shutdownFunc,
		errorHandlerImpl: errorHandler,
		logLevel:         log.INFO,
		shutdownDelay:    ShutdownDelay * time.Second,
	}

	if app.mux == nil {
		return nil, ErrInvalidConfiguration
	}

	// add the default endpoints

	// the root dir for the service's config
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	app.root = dir

	return app, nil
}
