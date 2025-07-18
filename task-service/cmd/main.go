package main

import (
	"github.com/cheeseong2001/task-service/handlers"
	"github.com/cheeseong2001/task-service/middleware"
	"github.com/cheeseong2001/task-service/repository"
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

	r.Use(middleware.JWTValidateMiddleware())

	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks", handlers.GetTask)
	r.GET("/tasks/:task_id", handlers.GetTaskByID)
	r.DELETE("/tasks/:task_id/delete", handlers.DeleteTask)
	r.Run(":8080")
}
