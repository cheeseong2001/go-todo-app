package models

type Task struct {
	TaskID      int    `json:"task_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	TaskID int `json:"task_id"`
}

type DeleteTaskRequest struct {
	TaskID int `json:"task_id"`
}

type DeleteTaskResponse struct {
	TaskID  int  `json:"task_id"`
	Success bool `json:"success"`
}
