package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/txsvc/apikit"
)

func TestNewClient(t *testing.T) {

	cl := NewClient(nil)
	assert.NotNil(t, cl)
}

func TestClientGET(t *testing.T) {

	cl := NewClient(nil)
	assert.NotNil(t, cl)

	// http setup
	svc, err := apikit.New(setup, shutdown)
	assert.NotNil(t, svc)
	assert.NoError(t, err)

	go func() {
		svc.Listen("")
	}()

	fmt.Println("waiting ...")
	time.Sleep(2 * time.Second)

	// the actual tests ...

	status, err := cl.GET("/test", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	status, err = cl.GET("/retry", nil)
	assert.Error(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, status)

	// http tear-down
	fmt.Println("waiting ...")
	time.Sleep(2 * time.Second)

	fmt.Println("stopping ...")
	// do NOT call svc.Stop() as it somehow does an os.Exit(1) thingie ...
}

func setup() *echo.Echo {
	e := echo.New()
	e.GET("/test", DefaultEndpoint)
	e.GET("/retry", testRetryEndpoint)

	return e
}

func shutdown(ctx context.Context, a *apikit.App) error {
	fmt.Println("shutting down ...")
	return nil
}

func testRetryEndpoint(c echo.Context) error {
	fmt.Println("delaying ...")
	time.Sleep(2 * time.Second)

	return StandardResponse(c, http.StatusServiceUnavailable, nil)
}
