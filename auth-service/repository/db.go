package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/cheeseong2001/auth-service/utils"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "user=postgres password=example dbname=auth-db host=auth-db port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	schema := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	DB = db
	return nil
}

func BootstrapAdminUser() error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1 AND role = 'admin'`, adminEmail).Scan(&count)
	if err != nil {
		return fmt.Errorf("error getting admin count in database: %w", err)
	}

	if count == 0 {
		hashedAdminPassword, err := utils.HashPassword(adminPassword)
		if err != nil {
			return fmt.Errorf("error hashing admin password: %w", err)
		}
		DB.Exec(`INSERT INTO users (email, password, role) VALUES ($1, $2, 'admin')`, adminEmail, hashedAdminPassword)
	}

	return nil
}
