package api

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"os"
)

//	@title						BrRSS
//	@version					1.0
//	@description				HTML to RSS Bridge
//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							query
//	@name						api-key

func EchoAPI() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	apiKey, apiKeySet := os.LookupEnv("API_KEY")
	apiUsername, apiUserSet := os.LookupEnv("API_USERNAME")
	apiPassword, apiPassSet := os.LookupEnv("API_PASSWORD")

	middlewares := []echo.MiddlewareFunc{
		middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "query:api-key",
			Skipper: func(ctx echo.Context) bool {
				return !apiKeySet
			},
			Validator: func(key string, ctx echo.Context) (bool, error) {
				return isValidKey([]byte(apiKey), []byte(key)), nil
			},
		}),
		middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Skipper: func(ctx echo.Context) bool {
				return !apiUserSet || !apiPassSet
			},
			Validator: func(username, password string, ctx echo.Context) (bool, error) {
				validUser := subtle.ConstantTimeCompare([]byte(username), []byte(apiUsername)) == 1
				validPass := isValidKey([]byte(apiPassword), []byte(password))

				return validUser && validPass, nil
			},
		}),
	}

	v1 := e.Group("/v1", middlewares...)

	v1.GET("/feed/:format", V1GetFeed)

	return e
}

func isValidKey(hashOrKey, key []byte) bool {
	_, err := bcrypt.Cost(hashOrKey)

	if err == nil {
		return bcrypt.CompareHashAndPassword(hashOrKey, key) == nil
	} else {
		return subtle.ConstantTimeCompare(hashOrKey, key) == 1
	}
}
