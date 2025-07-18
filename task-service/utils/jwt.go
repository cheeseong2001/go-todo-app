package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my-totally-secret-key") // need to replace with env var in prod env

func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, err
	}

	return int(userIDFloat), nil
}
