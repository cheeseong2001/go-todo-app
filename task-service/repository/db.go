package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "user=postgres password=example dbname=task-db host=task-db port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	schema := `CREATE TABLE IF NOT EXISTS tasks (
		task_id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT false
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	DB = db
	return nil
}
