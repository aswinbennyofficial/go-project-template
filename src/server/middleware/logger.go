package middlewares

import(
	"net/http"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"time"

)

// ZerologRequestLogger is a custom middleware that logs HTTP requests using zerolog
func ZerologRequestLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				logger.Info().
					Str("method", r.Method).
					Str("url", r.URL.String()).
					Int("status", ww.Status()).
					Dur("duration", time.Since(start)).
					Msg("handled request")
			}()
			next.ServeHTTP(ww, r)
		})
	}
}
