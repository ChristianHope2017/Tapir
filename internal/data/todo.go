package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/ChristianHope2017/di/internal/validator"
)

type Todo struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Task      string    `json:"task"`
}

func ValidateTodo(v *validator.Validator, todo *Todo) {
	v.Check(validator.NotBlank(todo.Title), "title", "must be provided")
	v.Check(validator.MaxLength(todo.Title, 50), "title", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(todo.Task), "task", "must be provided")
	v.Check(validator.MaxLength(todo.Task, 500), "task", "must not be more than 500 bytes long")
}

type TodoModel struct {
	DB *sql.DB
}

func (m *TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todo (title, task)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		todo.Title,
		todo.Task,
	).Scan(&todo.ID, &todo.CreatedAt)
}

func (m *TodoModel) Getall() ([]*Todo, error) {
	query := `SELECT id, created_at, title, task FROM todo ORDER BY created_at DESC`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.CreatedAt, &todo.Title, &todo.Task); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
