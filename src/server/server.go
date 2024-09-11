package server

import (
	"fmt"
	"net/http"

	"myapp/src/server/middleware"
	"myapp/src/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
	
)

type Server struct {
    App *utils.App
    Logger zerolog.Logger
}

func NewServer(app *utils.App) *Server {
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
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    r.Use(jwtauth.Verifier(tokenAuth))
    r.Use(jwtauth.Authenticator(tokenAuth))


    r.Get("/", s.HomeHandler)
    // Add more routes as needed

    addr := fmt.Sprintf(":%d", s.App.Config.App.Port)
    s.App.Logger.Info().Msgf("Starting server on %s", addr)
    http.ListenAndServe(addr, r)
}