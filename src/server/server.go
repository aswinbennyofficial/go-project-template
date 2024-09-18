package server

import (
	"fmt"
	"myapp/src/config"
	"myapp/src/server/handlers"
	"myapp/src/server/middleware"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type Server struct {
	App    *config.App
	Logger zerolog.Logger
}

func NewServer(app *config.App) *Server {
	return &Server{
		App:    app,
		Logger: app.Logger,
	}
}

func (s *Server) Start() {
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middlewares.ZerologLogger(s.Logger))
	e.Use(middleware.Recover())

	// JWT middleware
	jwtConfig := echojwt.Config{
		SigningKey: []byte(s.App.Config.Auth.JWTSecret),
	}

	// Routes
	v1 := e.Group("/api/v1")
	v1.Use(echojwt.WithConfig(jwtConfig))
	v1.GET("/home", handlers.HomeHandler(s.App, s.Logger))

	// Start server
	addr := fmt.Sprintf(":%d", s.App.Config.App.Port)
	s.App.Logger.Info().Msgf("Starting server on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}