package apikit

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/txsvc/apikit/config"
	"github.com/ziflex/lecho/v3"
)

const (
	// default ports to listen on
	PORT_DEFAULT     = "8080"
	PORT_DEFAULT_TLS = "433"

	// ShutdownDelay is the time to wait for all request, go-routines etc complete
	ShutdownDelay = 30 // seconds
)

type (
	// SetupFunc creates a new, fully configured mux
	SetupFunc func() *echo.Echo
	// ShutdownFunc is called before the server stops
	ShutdownFunc func(context.Context, *App) error

	// app holds all configs for the listener
	App struct {
		mux      *echo.Echo
		shutdown ShutdownFunc

		// other settings
		logLevel      log.Lvl
		shutdownDelay time.Duration
		root          string
	}
)

// New creates a new service listener instance and configures it with sensible defaults.
//
// The following ENV variables are supported:
// - PORT: default 8080
// - CONFIG_LOCATION: default ./.config
// - LOG_LEVEL: default INFO
//
// - force_ssl: default false
// - secret_key_base:
// - public_file_server: default false
func New(setupFunc SetupFunc, shutdownFunc ShutdownFunc) (*App, error) {
	if setupFunc == nil || shutdownFunc == nil {
		return nil, config.ErrInvalidConfiguration
	}

	app := &App{
		mux:           setupFunc(),
		shutdown:      shutdownFunc,
		logLevel:      log.INFO,
		shutdownDelay: ShutdownDelay * time.Second,
	}

	if app.mux == nil {
		return nil, config.ErrInvalidConfiguration
	}

	// no greetings
	app.mux.HideBanner = true

	// add a logger and middleware
	logger := lecho.New(os.Stdout)
	app.mux.Logger.SetLevel(app.logLevel)
	app.mux.Logger = logger

	app.mux.Use(middleware.RequestID())
	app.mux.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))

	// add a default error handler
	// app.mux.HTTPErrorHandler = ...

	// the root dir for the config
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	app.root = dir

	return app, nil
}

func (a *App) Stop() {
	// FIXME: this does not work !
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownDelay)
	defer cancel()

	// FIXME: which one comes first ? framwork or app shutdown ?

	// call the implementation specific shoutdown code to clean-up
	if err := a.shutdown(ctx, a); err != nil {
		a.mux.Logger.Fatal(err)
	}

	// shutdown of the framework
	if err := a.mux.Shutdown(ctx); err != nil {
		a.mux.Logger.Fatal(err)
	}
}
