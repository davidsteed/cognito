package main

import (
	"net/http"

	Openapi "github.com/davidsteed/cognito/api"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetLocation(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, []Openapi.Location{})
}
