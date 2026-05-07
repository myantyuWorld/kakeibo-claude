package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
}

func NewServer(origins []string) *Server {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
	}))
	e.HTTPErrorHandler = errorHandler
	e.Validator = newValidator()
	e.HideBanner = true

	return &Server{Echo: e}
}

func (s *Server) Run(addr string) error {
	return s.Start(":" + addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}
