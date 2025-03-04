package handlers

import (
	"context"
	"fmt"
	"myapp/src/config"
	"myapp/src/utils"
	"net/http"

	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

// PingHandler pings all the databases and dependencies to check if they are up
func PingHandler(app *config.App, logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Postgres  string `json:"postgres"`
			Redis     string `json:"redis"`
			Cassandra string `json:"cassandra"`
			RabbitMQ string `json:"rabbitmq"`
		}

		resp := response{}

		// Ping Postgres
		if err := app.Postgres.Ping(context.Background()); err != nil {
			logger.Error().Err(err).Msg("Failed to ping postgres")
			resp.Postgres = fmt.Sprintf("err: %s", err)
		} else {
			resp.Postgres = "ok"
		}

		// Ping Redis
		if err := app.Redis.Ping(context.Background()).Err(); err != nil {
			logger.Error().Err(err).Msg("Failed to ping redis")
			resp.Redis = fmt.Sprintf("err: %s", err)
		} else {
			resp.Redis = "ok"
		}

		// Ping Cassandra
		if err := app.Cassandra.Query("SELECT release_version FROM system.local").Exec(); err != nil {
			logger.Error().Err(err).Msg("Failed to ping cassandra")
			resp.Cassandra = fmt.Sprintf("err: %s", err)
		} else {
			resp.Cassandra = "ok"
		}

		// Ping RabbitMQ
		if _, err := app.RabbitMQ.Channel(); err != nil {
			logger.Error().Err(err).Msg("Failed to ping rabbitmq")
			resp.RabbitMQ = fmt.Sprintf("err: %s", err)
		} else {
			resp.RabbitMQ = "ok"
		}

		// If any dependency is down, return HTTP 500
		if resp.Postgres != "ok" || resp.Redis != "ok" || resp.Cassandra != "ok" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			render.Status(r, http.StatusOK)
		}
		render.JSON(w, r, resp)
	}
}

// MetricsHandler handles the /metrics route returns the prometheus metrics
func MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}

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
