package main

import (
	"github.com/cheeseong2001/auth-service/handlers"
	"github.com/cheeseong2001/auth-service/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	err := repository.InitDB()
	if err != nil {
		panic(err)
	}
	err = repository.DB.Ping()
	if err != nil {
		panic(err)
	}
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/validate", handlers.ValidateToken)
	}

	r.Run(":8080")
}
