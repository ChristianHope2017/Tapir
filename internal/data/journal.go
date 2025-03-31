package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/ChristianHope2017/di/internal/validator"
)

type Journal struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

func ValidateJournal(v *validator.Validator, journal *Journal) {
	v.Check(validator.NotBlank(journal.Title), "title", "must be provided")
	v.Check(validator.MaxLength(journal.Title, 50), "title", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(journal.Content), "content", "must be provided")
	v.Check(validator.MaxLength(journal.Content, 500), "content", "must not be more than 500 bytes long")
}

type JournalModel struct {
	DB *sql.DB
}

func (m *JournalModel) Insert(journal *Journal) error {
	query := `
		INSERT INTO journal (title, content)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		journal.Title,
		journal.Content,
	).Scan(&journal.ID, &journal.CreatedAt)
}

func (m *JournalModel) GetAll() ([]*Journal, error) {
	query := `SELECT id, created_at, title, content FROM journal ORDER BY created_at DESC`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journals := []*Journal{}
	for rows.Next() {
		journal := &Journal{}
		err = rows.Scan(&journal.ID, &journal.CreatedAt, &journal.Title, &journal.Content)
		if err != nil {
			return nil, err
		}
		journals = append(journals, journal)
	}

	return journals, nil
}
