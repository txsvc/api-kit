package apikit

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"

	"github.com/txsvc/apikit/config"
)

const (
	// default ports to listen on
	PORT_DEFAULT     = "8080"
	PORT_DEFAULT_TLS = "433"

	// ShutdownDelay is the time to wait for all request, go-routines etc complete
	ShutdownDelay = 30 // seconds
)

type (
	// SetupFunc creates a new, fully configured router
	SetupFunc func() *echo.Echo
	// ShutdownFunc is called before the app stops
	ShutdownFunc func(context.Context, *App) error

	// app holds all configs for the listener
	App struct {
		svc *echo.Echo

		shutdown      ShutdownFunc
		shutdownDelay time.Duration

		// other settings
		logLevel log.Lvl
		root     string
	}
)

// New creates a new service listener instance and configures it with sensible defaults.
func New(setupFunc SetupFunc, shutdownFunc ShutdownFunc) (*App, error) {
	if setupFunc == nil || shutdownFunc == nil {
		return nil, config.ErrInvalidConfiguration
	}

	app := &App{
		svc:           setupFunc(),
		shutdown:      shutdownFunc,
		logLevel:      log.INFO,
		shutdownDelay: ShutdownDelay * time.Second,
	}

	if app.svc == nil {
		return nil, config.ErrInvalidConfiguration
	}

	// no greetings
	app.svc.HideBanner = true

	// add a logger and middleware
	logger := lecho.New(os.Stdout)
	app.svc.Logger = logger
	app.svc.Logger.SetLevel(app.logLevel)

	// adding logging related middleware
	app.svc.Use(middleware.RequestID())
	app.svc.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))

	// FIXME: add a default error handler
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
		a.svc.Logger.Fatal(err)
	}

	// shutdown of the framework
	if err := a.svc.Shutdown(ctx); err != nil {
		a.svc.Logger.Fatal(err)
	}
}
