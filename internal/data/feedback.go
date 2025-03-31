package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/ChristianHope2017/di/internal/validator"
)

type Feedback struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Fullname  string    `json:"fullname"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Email     string    `json:"email"`
}

func ValidateFeedback(v *validator.Validator, feedback *Feedback) {
	v.Check(validator.NotBlank(feedback.Fullname), "fullname", "must be provided")
	v.Check(validator.MaxLength(feedback.Fullname, 50), "fullname", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(feedback.Subject), "subject", "must be provided")
	v.Check(validator.MaxLength(feedback.Subject, 50), "subject", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(feedback.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(feedback.Email), "email", "invalid email address")
	v.Check(validator.MaxLength(feedback.Email, 100), "email", "must not be more than 100 bytes long")
	v.Check(validator.NotBlank(feedback.Message), "message", "must be provided")
	v.Check(validator.MaxLength(feedback.Message, 500), "message", "must not be more than 500 bytes long")
}

type FeedbackModel struct {
	DB *sql.DB
}

func (m *FeedbackModel) Insert(feedback *Feedback) error {
	query := `
		INSERT INTO feedback (fullname, subject, message, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		feedback.Fullname,
		feedback.Subject,
		feedback.Message,
		feedback.Email,
	).Scan(&feedback.ID, &feedback.CreatedAt)
}

func (m *FeedbackModel) GetAll() ([]*Feedback, error) {
	query := `SELECT id, created_at, fullname, subject, message, email FROM feedback ORDER BY created_at DESC`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := []*Feedback{}
	for rows.Next() {
		feedback := &Feedback{}
		err = rows.Scan(&feedback.ID, &feedback.CreatedAt, &feedback.Fullname, &feedback.Subject, &feedback.Message, &feedback.Email)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}
