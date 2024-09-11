package server

import (
	"fmt"
	"myapp/src/utils"
	"net/http"
)

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Use the utility function to get the user_id from the JWT token
	userID, err := utils.ExtractClaim(r.Context(), "user_id")
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to extract user_id from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Log and respond with a message
	s.Logger.Info().Str("user_id", userID).Msg("Home page accessed")
	response := fmt.Sprintf("Welcome to MyApp, %s!", userID)
	w.Write([]byte(response))
}