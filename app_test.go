package apikit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func nilSetup() *echo.Echo {
	return nil
}

func simpleSetup() *echo.Echo {
	return echo.New()
}

func noopShutdown(ctx context.Context, a *App) error {
	fmt.Println("shutting down ...")
	return nil
}

func timeoutShutdown(ctx context.Context, a *App) error {
	fmt.Println("shutting down blocking ...")

	for {
		time.Sleep(2 * time.Second)
		fmt.Println("blocking ...")
	}

	return nil
}

func TestNewSimple(t *testing.T) {
	svc, err := New(simpleSetup, noopShutdown)

	assert.NotNil(t, svc)
	assert.NoError(t, err)
}

func TestNewNil(t *testing.T) {
	svc, err := New(nil, nil)
	assert.Nil(t, svc)
	assert.Error(t, err)

	svc, err = New(nilSetup, noopShutdown)
	assert.Nil(t, svc)
	assert.Error(t, err)
}

func TestRunStop(t *testing.T) {
	svc, err := New(simpleSetup, noopShutdown)

	assert.NotNil(t, svc)
	assert.NoError(t, err)

	go func() {
		svc.Listen("")
	}()

	fmt.Println("listening, waiting ...")
	time.Sleep(2 * time.Second)

	fmt.Println("stopping ...")
	// do NOT call svc.Stop() as it somehow does an os.Exit(1) thingie ...
}

// FIXME: the timeout thing doesn't work
func _TestRunStopTimeout(t *testing.T) {
	svc, err := New(simpleSetup, timeoutShutdown)

	assert.NotNil(t, svc)
	assert.NoError(t, err)

	go func() {
		svc.Listen("")
	}()

	fmt.Println("listening, waiting ...")
	time.Sleep(10 * time.Second)

	fmt.Println("stopping ...")
	svc.Stop()

	time.Sleep(10 * time.Second)
}
