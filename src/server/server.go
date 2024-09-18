package server

import (
	"fmt"
	"myapp/src/config"
	"myapp/src/server/handlers"
	middlewares "myapp/src/server/middleware"

	"github.com/gin-gonic/gin"
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
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middlewares.ZerologLogger(s.Logger))

	// JWT middleware
	authMiddleware, err := middlewares.InitJWTAuth(s.App.Config.Auth.JWTSecret)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to initialize JWT middleware")
	}

	// Routes
	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		v1.GET("/home", handlers.HomeHandler(s.App, s.Logger))
	}

	// Start server
	addr := fmt.Sprintf(":%d", s.App.Config.App.Port)
	s.Logger.Info().Msgf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}