package apikit

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v3"

	"github.com/txsvc/stdlib/v2"
	"github.com/txsvc/stdlib/v2/stdlibx/stringsx"
)

func (a *App) Listen(addr string) {
	// setup shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		a.Stop()
	}()

	// add the central error handler
	if a.errorHandlerImpl != nil {
		a.mux.HTTPErrorHandler = a.errorHandlerImpl
	}

	a.mux.HideBanner = true

	// add a logger and middleware
	logger := lecho.New(os.Stdout)
	a.mux.Logger.SetLevel(a.logLevel)
	a.mux.Logger = logger

	a.mux.Use(middleware.RequestID())
	a.mux.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))

	port := fmt.Sprintf(":%s", stringsx.TakeOne(stdlib.GetString("PORT", addr), "8080"))
	a.mux.Logger.Fatal(a.mux.Start(port))
}

func (a *App) ListenAutoTLS(addr string) {

}

func (a *App) Stop() {
	// all the implementation specific shoutdown code to clean-up
	a.shutdown(a)

	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownDelay)
	defer cancel()
	if err := a.mux.Shutdown(ctx); err != nil {
		a.mux.Logger.Fatal(err)
	}
}
