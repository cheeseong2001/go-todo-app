package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthResponse struct {
	Valid  bool   `json:"valid"`
	UserID int    `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
	Error  string `json:"error,omitempty"`
}

func JWTValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid auth header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		body := map[string]string{"token": tokenString}
		jsonBody, _ := json.Marshal(body)

		resp, err := http.Post("http://auth-service:8080/auth/validate", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil || resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token validation failed",
			})
			return
		}
		defer resp.Body.Close()

		var authResponse AuthResponse

		if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse auth response",
			})
			return
		}

		c.Set("user_id", authResponse.UserID)
		c.Set("role", authResponse.Role)
		c.Next()
	}
}
