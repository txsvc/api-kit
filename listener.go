package apikit

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/txsvc/stdlib/v2"
	"github.com/txsvc/stdlib/v2/stdlibx/stringsx"
)

func (a *App) listen(addr, certFile, keyFile string, useTLS bool) {
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

	if useTLS {
		port := fmt.Sprintf(":%s", stringsx.TakeOne(stdlib.GetString("PORT", addr), PORT_DEFAULT_TLS))
		certDir := fmt.Sprintf("%s/.cert", a.root)

		autoTLSManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
			Cache: autocert.DirCache(certDir),
			//HostPolicy: autocert.HostWhitelist("<DOMAIN>"),
		}
		tlsc := tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		}
		s := http.Server{
			Addr:      port,
			Handler:   a.mux, // set Echo as handler
			TLSConfig: &tlsc,
		}
		if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	} else {
		// simply startup without TLS
		port := fmt.Sprintf(":%s", stringsx.TakeOne(stdlib.GetString("PORT", addr), PORT_DEFAULT))
		log.Fatal(a.mux.Start(port))
	}
}

func (a *App) Listen(addr string) {
	a.listen(addr, "", "", false)
}

func (a *App) ListenAutoTLS(addr string) {
	a.listen(addr, "", "", true)
}

func (a *App) ListenTLS(addr, certFile, keyFile string) {
	a.listen(addr, certFile, keyFile, true)
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
