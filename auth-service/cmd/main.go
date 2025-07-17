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

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/validate", handlers.ValidateToken)
	r.Run(":8080")
}
