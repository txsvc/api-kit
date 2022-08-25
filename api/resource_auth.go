package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	// auth routes
	InitRoute   = "/auth/init"
	LoginRoute  = "/auth/login"
	LogoutRoute = "/auth/logout"
)

func WithAuthEndpoints(e *echo.Echo) *echo.Echo {
	// grouped under /a/v1
	apiGroup := e.Group(NamespacePrefix)

	// add the routes
	apiGroup.POST(InitRoute, InitEndpoint)
	apiGroup.POST(LoginRoute, LoginEndpoint)
	apiGroup.PUT(LogoutRoute, LogoutEndpoint)

	// done
	return e
}

func InitEndpoint(c echo.Context) error {
	return StandardResponse(c, http.StatusNotImplemented, nil)
}

func LoginEndpoint(c echo.Context) error {
	return StandardResponse(c, http.StatusNotImplemented, nil)
}

func LogoutEndpoint(c echo.Context) error {
	return StandardResponse(c, http.StatusNotImplemented, nil)
}
