package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func ValidateJWT(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", err
	}

	// token has expired
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return 0, "", errors.New("token expired")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", err
	}

	return int(userIDFloat), role, nil
}
