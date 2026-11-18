package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mushibana/glacier/backend/internal/domain"
)

type tasksRow struct {
	ID       int64
	ParentID sql.NullInt64
	OrderNum int64

	Name        string
	Description sql.NullString
	StreamID    sql.NullInt64
	StartAt     sql.NullString
	DeadlineAt  sql.NullString
	TimeLimit   int64
	IsOnetime   bool
	Repeatable  bool
	Finished    bool

	CreatedAt sql.NullString
	UpdatedAt sql.NullString
}

func toDomain(r tasksRow) (domain.Task, error) {
	t := domain.Task{
		ID:          r.ID,
		ParentID:    fromNullableInt64ToPointer(r.ParentID),
		OrderNum:    r.OrderNum,
		Name:        r.Name,
		Description: r.Description.String,
		StreamID:    fromNullableInt64ToPointer(r.StreamID),
		TimeLimit:   r.TimeLimit,
		IsOnetime:   r.IsOnetime,
		Repeatable:  r.Repeatable,
		Finished:    r.Finished,
	}

	var err error

	t.StartAt, err = fromNullableTime(r.StartAt, time.RFC3339)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parse StartAt: %w", err)
	}
	t.DeadlineAt, err = fromNullableTime(r.DeadlineAt, time.RFC3339)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parse DeadlineAt: %w", err)
	}
	t.CreatedAt, err = fromNullableTime(r.CreatedAt, time.RFC3339)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parse CreatedAt: %w", err)
	}
	t.UpdatedAt, err = fromNullableTime(r.UpdatedAt, time.RFC3339)
	if err != nil {
		return domain.Task{}, fmt.Errorf("parse UpdatedAt: %w", err)
	}

	return t, nil
}

func fromDomain(t domain.Task) (tasksRow, error) {
	r := tasksRow{
		ID:          t.ID,
		ParentID:    toNullInt64FromPointer(t.ParentID),
		OrderNum:    t.OrderNum,
		Name:        t.Name,
		Description: toNullString(t.Description),
		StreamID:    toNullInt64FromPointer(t.StreamID),
		StartAt:     toNullTime(t.StartAt, time.RFC3339),
		DeadlineAt:  toNullTime(t.DeadlineAt, time.RFC3339),
		TimeLimit:   t.TimeLimit,
		IsOnetime:   t.IsOnetime,
		Repeatable:  t.Repeatable,
		Finished:    t.Finished,
		CreatedAt:   toNullTime(t.CreatedAt, time.RFC3339),
		UpdatedAt:   toNullTime(t.UpdatedAt, time.RFC3339),
	}

	return r, nil
}
