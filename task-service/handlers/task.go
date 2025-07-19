package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/cheeseong2001/task-service/models"
	"github.com/cheeseong2001/task-service/repository"
	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var createTaskReq models.CreateTaskRequest

	if err := c.ShouldBindJSON(&createTaskReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "user not found in context",
		})
		return
	}

	var taskID int
	repository.DB.QueryRow(`INSERT INTO tasks (user_id, title, description, completed) 
	VALUES ($1, $2, $3, $4) RETURNING task_id;`, userID, createTaskReq.Title, createTaskReq.Description, false).Scan(&taskID)

	c.IndentedJSON(http.StatusCreated, models.CreateTaskResponse{
		TaskID: taskID,
	})
}

func GetTask(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "role not found in context",
		})
		return
	}

	var query string
	var rows *sql.Rows
	var err error

	if role == "admin" { // give admin power to list all tasks
		query = `SELECT task_id, user_id, title, description, completed
		FROM tasks;`
		rows, err = repository.DB.Query(query)

	} else { // give users power to list their own tasks
		userID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "user not found in context",
			})
			return
		}
		query = `SELECT task_id, user_id, title, description, completed
		FROM tasks
		WHERE user_id = $1;`
		rows, err = repository.DB.Query(query, userID)
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.TaskID, &task.UserID, &task.Title, &task.Description, &task.Completed); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "scan errror",
			})
			return
		}

		tasks = append(tasks, task)
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "user not found in context",
		})
		return
	}

	var task models.Task
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task_id",
		})
		return
	}
	task.TaskID = taskID

	err = repository.DB.QueryRow(`
		SELECT title, description, completed
		FROM tasks
		WHERE user_id = $1 AND task_id = $2;`, userID, taskID).Scan(&task.Title, &task.Description, &task.Completed)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "task not found",
			})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "user not found in context",
		})
		return
	}

	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid task_id",
		})
		return
	}

	result, err := repository.DB.Exec(`DELETE FROM tasks WHERE task_id = $1 AND user_id = $2`, taskID, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete task",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "task_id not found",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, models.DeleteTaskResponse{
		TaskID:  taskID,
		Success: true,
	})
}

func UpdateTaskComplete(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "role not found in context",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "user not found in context",
		})
		return
	}

	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid task_id",
		})
		return
	}

	var updateTaskCompletedReq models.UpdateTaskCompletedRequest

	if err = c.ShouldBindJSON(&updateTaskCompletedReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var query string
	var result sql.Result

	if role == "admin" {
		query = `UPDATE tasks SET completed = $1 WHERE task_id = $2`
		result, err = repository.DB.Exec(query, updateTaskCompletedReq.Completed, taskID)
	} else {
		query = `UPDATE tasks SET completed = $1 WHERE user_id = $2 AND task_id = $3`
		result, err = repository.DB.Exec(query, updateTaskCompletedReq.Completed, userID, taskID)
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update task",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "task_id not found",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, models.UpdateTaskCompletedResponse{
		TaskID:    taskID,
		Completed: updateTaskCompletedReq.Completed,
		Success:   true,
	})
}
