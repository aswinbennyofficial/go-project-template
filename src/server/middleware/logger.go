package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func ZerologLogger(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			logger.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("remote_ip", c.RealIP()).
				Int("status", res.Status).
				Msg("Request")

			return next(c)
		}
	}
}