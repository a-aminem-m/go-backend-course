package main

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres() (*sql.DB, error) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		dsn = "postgres://app:app@localhost:5432/app?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec(
		`INSERT INTO users (id, login, password)
 		 VALUES ($1, $2, $3)
 		 ON CONFLICT (login) DO NOTHING`,
		user.ID,
		user.Login,
		user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(db *sql.DB, login string) (User, bool) {
	var user User

	err := db.QueryRow(
		`SELECT id, login, password
		 FROM users
		 WHERE login = $1`,
		login,
	).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return User{}, false
	}

	if err != nil {
		return User{}, false
	}

	return user, true
}

func CreateTask(db *sql.DB, task Task) error {
	_, err := db.Exec(
		`INSERT INTO tasks (id, status, result, translator, code)
		 VALUES ($1, $2, $3, $4, $5)`,
		task.ID,
		task.Status,
		task.Result,
		task.Translator,
		task.Code,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetTask(db *sql.DB, id string) (Task, bool) {
	var task Task

	err := db.QueryRow(
		`SELECT id, status, result, translator, code
		 FROM tasks
		 WHERE id = $1`,
		id,
	).Scan(
		&task.ID,
		&task.Status,
		&task.Result,
		&task.Translator,
		&task.Code,
	)

	if err == sql.ErrNoRows {
		return Task{}, false
	}

	if err != nil {
		return Task{}, false
	}

	return task, true
}

func UpdateTask(db *sql.DB, task Task) error {
	_, err := db.Exec(
		`UPDATE tasks
		 SET status = $2, result = $3, translator = $4, code = $5
		 WHERE id = $1`,
		task.ID,
		task.Status,
		task.Result,
		task.Translator,
		task.Code,
	)
	if err != nil {
		return err
	}

	return nil
}
