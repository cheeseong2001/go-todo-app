package middleware

import (
	"net/http"
	"strings"

	"github.com/cheeseong2001/task-service/utils"
	"github.com/gin-gonic/gin"
)

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
		userID, role, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}
