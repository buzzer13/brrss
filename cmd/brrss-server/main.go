package main

import (
	"github.com/buzzer13/brrss/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/buzzer13/brrss/docs"
)

func main() {
	a := &api.API{}
	e := a.New()

	e.Use(middleware.Logger())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(307, "/swagger/index.html")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
