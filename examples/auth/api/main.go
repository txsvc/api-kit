package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/txsvc/apikit"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/internal"
)

func init() {
	// initialize the config provider
	config.InitConfigProvider(config.NewLocalConfigProvider())

	// create a default configuration for the service (if none exists)
	path := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		// create credentials and keys
		cfg, err := internal.InitSettings(config.Name(), config.Name())
		if err != nil {
			log.Fatal(err)
		}

		// defaults from this config provider
		def := config.GetDefaultSettings()

		// copy the credentials and api keys
		def.Credentials = cfg.Credentials
		def.APIKey = cfg.APIKey
		def.Scopes = append(def.Scopes, internal.ScopeApiAdmin)

		// save the new configuration
		def.WriteToFile(path)
	}

	// initialize the credentials store
	root := filepath.Join(config.ResolveConfigLocation(), "cred")
	internal.FlushAuthorizations(root)
}

func main() {
	// example of using a cmd line flag for configuration
	useTLS := flag.Bool("tls", false, "use TLS endpoint termination")
	flag.Parse()

	svc, err := apikit.New(setup, shutdown)
	if err != nil {
		log.Fatal(err)
	}

	if !*useTLS {
		// start listening on default port 8080
		svc.Listen("")
	} else {
		// start listening on default port 443
		svc.ListenAutoTLS("")
	}
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
	_, err := internal.CheckAuthorization(ctx, c, internal.ScopeApiRead)
	if err != nil {
		return apikit.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	resp := apikit.StatusObject{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("version: %s", config.VersionString()),
	}

	return apikit.StandardResponse(c, http.StatusOK, resp)
}
