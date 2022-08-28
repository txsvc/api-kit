package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/txsvc/apikit"
	"github.com/txsvc/apikit/api"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/helpers"
	"github.com/txsvc/apikit/internal/auth"
)

const (
	// MajorVersion of the API
	majorVersion = 0
	// MinorVersion of the API
	minorVersion = 1
	// FixVersion of the API
	fixVersion = 0
)

type (
	appConfig struct {
		root string // the fully qualified path to the conf dir
		info *config.Info
	}
)

// FIXME: implement in memory certstore

func init() {
	// initialize the config provider
	config.InitConfigProvider(NewAppEngineConfigProvider())

	// create a default configuration for the service (if none exists)
	path := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		// create credentials and keys with defaults from this config provider
		cfg := config.GetDefaultSettings()

		// save the new configuration
		helpers.WriteSettingsToFile(cfg, path)
	}

	// initialize the credentials store
	root := filepath.Join(config.ResolveConfigLocation(), "cred")
	auth.FlushAuthorizations(root)
}

func main() {

	svc, err := apikit.New(setup, shutdown)
	if err != nil {
		log.Fatal(err)
	}

	// Do not use AutoTLS here as TLS termination is handled by App Engine.
	// Do not change default port 8080 !
	svc.Listen("")
}

func setup() *echo.Echo {
	// create a new router instance
	e := echo.New()

	// add and configure any middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// add your endpoints here
	e.GET("/ping", pingEndpoint)

	// done
	return e
}

func shutdown(ctx context.Context, a *apikit.App) error {
	// TODO: implement your own stuff here
	return nil
}

// pingEndpoint returns http.StatusOK and the version string
func pingEndpoint(c echo.Context) error {
	ctx := context.Background()

	// this endpoint needs at minimum an "api:read" scope
	_, err := auth.CheckAuthorization(ctx, c, auth.ScopeApiRead)
	if err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err, "")
	}

	resp := api.StatusObject{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("version: %s", config.AppInfo().VersionString()),
	}

	return api.StandardResponse(c, http.StatusOK, resp)
}

func NewAppEngineConfigProvider() interface{} {
	info := config.NewAppInfo(
		"appengine kit",
		"aek",
		"Copyright 2022, transformative.services, https://txs.vc",
		"about appengine kit",
		majorVersion,
		minorVersion,
		fixVersion,
	)

	return &appConfig{
		info: &info,
	}
}
