package main

import (
	"github.com/buzzer13/brrss/api"
	"github.com/buzzer13/gosls/do"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var a = api.API{}

func Main(evm do.FuncEventMap) (do.FuncResponseMap, error) {
	res := do.FuncResponseWriter{}
	evt, err := evm.Event()

	if err != nil {
		return (&do.FuncResponse{
			Body:       "event parse error - " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}).Map(), err
	}

	req, err := evt.Request()

	if err != nil {
		return (&do.FuncResponse{
			Body:       "request parse error - " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}).Map(), err
	}

	e := a.NewOnce(func(a *api.API) {
		a.Echo.Use(middleware.Logger())
	})

	e.ServeHTTP(&res, req)

	return res.GetFuncResponse().Map(), nil
}
