package domain

import "time"

type Task struct {
	ID       int64
	ParentID *int64
	OrderNum int64

	Name        string
	Description string
	StreamID    *int64
	StartAt     time.Time
	DeadlineAt  time.Time
	TimeLimit   int64
	IsOnetime   bool
	Repeatable  bool
	Finished    bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
