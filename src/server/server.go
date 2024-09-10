package server

import (
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "myapp/src/server/middleware"
    "myapp/src/utils"
)

type Server struct {
    App *utils.App
}

func NewServer(app *utils.App) *Server {
    return &Server{App: app}
}

func (s *Server) Start() {
    r := chi.NewRouter()

    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middlewares.JWTAuth(s.App.Config.Auth.JWTSecret))

    r.Get("/", s.HomeHandler)
    // Add more routes as needed

    addr := fmt.Sprintf(":%d", s.App.Config.App.Port)
    s.App.Logger.Info().Msgf("Starting server on %s", addr)
    http.ListenAndServe(addr, r)
}