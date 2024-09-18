package handlers

import (
	"fmt"
	"myapp/src/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	jwt "github.com/appleboy/gin-jwt/v2"
)

// HomeHandler handles the home route
func HomeHandler(app *config.App, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from JWT
		claims := jwt.ExtractClaims(c)
		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.Error().Msg("Failed to extract user_id from token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Log and respond with a message
		logger.Info().Str("user_id", userID).Msg("Home page accessed")

		// Create the response payload
		response := gin.H{
			"message": fmt.Sprintf("Welcome to MyApp, %s!", userID),
		}

		// Respond with JSON
		c.JSON(http.StatusOK, response)
	}
}