package server

import (
    "net/http"
)

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to MyApp!"))
}
