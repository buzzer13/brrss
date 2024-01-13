package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

// @title			BrRSS
// @version		1.0
// @description	HTML to RSS Bridge

func EchoAPI() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	apiKey, apiKeySet := os.LookupEnv("API_KEY")

	v1 := e.Group("/v1", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key",
		Skipper: func(ctx echo.Context) bool {
			return !apiKeySet
		},
		Validator: func(key string, ctx echo.Context) (bool, error) {
			return key == apiKey, nil
		},
	}))

	v1.GET("/feed/:format", V1GetFeed)

	return e
}
