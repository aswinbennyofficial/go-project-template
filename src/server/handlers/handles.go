package handlers

import (
	"fmt"
	"myapp/src/config"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func HomeHandler(app *config.App, logger zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract user from JWT
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Log and respond with a message
		logger.Info().Str("user_id", userID).Msg("Home page accessed")

		// Create the response payload
		response := map[string]string{
			"message": fmt.Sprintf("Welcome to MyApp, %s!", userID),
		}

		return c.JSON(http.StatusOK, response)
	}
}