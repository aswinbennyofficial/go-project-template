package middlewares

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt"
)

func JWTAuth(jwtSecret string) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Missing auth token", http.StatusUnauthorized)
                return
            }

            bearerToken := strings.Split(authHeader, " ")
            if len(bearerToken) != 2 {
                http.Error(w, "Invalid auth token", http.StatusUnauthorized)
                return
            }

            token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid auth token", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                http.Error(w, "Invalid auth token", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), "user", claims["user"])
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}