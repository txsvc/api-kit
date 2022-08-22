package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/txsvc/apikit"
)

func init() {

}

func setup() *echo.Echo {
	// create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// TODO add your own endpoints here
	e.GET("/", apikit.DefaultEndpoint)

	return e
}

func shutdown(a *apikit.App) {
	// TODO: implement your own stuff here
}

func main() {
	svc, err := apikit.New(setup, shutdown, nil)
	if err != nil {
		log.Fatal(err)
	}

	svc.Listen("")
}
