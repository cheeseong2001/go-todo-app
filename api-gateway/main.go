package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	authServiceURL = "http://auth-service:8080"
	taskServiceURL = "http://task-service:8080"
)

func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, err := url.Parse(target)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "bad URL",
			})
			return
		}

		reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.Default()

	r.Any("/auth/*proxyPath", reverseProxy(authServiceURL))
	r.Any("/tasks/*proxyPath", reverseProxy(taskServiceURL))

	r.Run(":8080")
}
