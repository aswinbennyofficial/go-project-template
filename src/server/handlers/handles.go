package handlers

import (
	"fmt"
	"myapp/src/config"
	"myapp/src/utils"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/go-chi/render"
)

// HomeHandler handles the home route
func HomeHandler(app *config.App, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT
		userID, err := utils.ExtractClaim(r.Context(), "user_id")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to extract user_id from token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Log and respond with a message
		logger.Info().Str("user_id", userID).Msg("Home page accessed")
		
		// Create the response payload
		response := map[string]string{
			"message": fmt.Sprintf("Welcome to MyApp, %s!", userID),
		}

		// Use Chi render to bind and respond with JSON
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)

	}
}
