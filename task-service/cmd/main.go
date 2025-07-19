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

	tasks := r.Group("/tasks")
	{
		tasks.POST("/", handlers.CreateTask)
		tasks.GET("/", handlers.GetTask)
		tasks.GET("/:task_id", handlers.GetTaskByID)
		tasks.DELETE("/:task_id", handlers.DeleteTask)
		tasks.PATCH("/:task_id/complete", handlers.UpdateTaskComplete)
	}

	r.Run(":8080")
}
