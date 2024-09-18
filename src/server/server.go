package server

import (
	"fmt"
	"net/http"

	"myapp/src/config"
	"myapp/src/server/handlers"
	"myapp/src/server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
)

type Server struct {
    App *config.App
    Logger zerolog.Logger
}

func NewServer(app *config.App) *Server {
    return &Server{
		App:    app,
		Logger: app.Logger,
	}
}

func (s *Server) Start() {
    r := chi.NewRouter()


    tokenAuth := middlewares.InitJWTAuth(s.App.Config.Auth.JWTSecret)

    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middlewares.ZerologRequestLogger(s.Logger))
    r.Use(middleware.Recoverer)
    

    r.Route("/api/v1", func(r chi.Router) {
        r.Use(jwtauth.Verifier(tokenAuth))
        r.Use(jwtauth.Authenticator(tokenAuth))
    
        r.Get("/home", handlers.HomeHandler(s.App, s.Logger)) 
    })

    addr := fmt.Sprintf(":%d", s.App.Config.App.Port)
    s.App.Logger.Info().Msgf("Starting server on %s", addr)
    http.ListenAndServe(addr, r)
}