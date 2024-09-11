package utils

import (
	"context"
	"errors"

	"github.com/go-chi/jwtauth/v5"

)



// ExtractClaim retrieves a specific claim (like "user_id") from the JWT token in the context
func ExtractClaim(ctx context.Context, key string) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", err
	}

	value, ok := claims[key].(string)
	if !ok {
		return "", errors.New("claim not found or not a string")
	}

	return value, nil
}