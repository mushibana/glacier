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

type UpdateTaskParams struct {
	ID          int64
	ParentID    Optional[*int64]
	Name        Optional[string]
	Description Optional[string]
	StreamID    Optional[*int64]
	StartAt     Optional[time.Time]
	DeadlineAt  Optional[time.Time]
	TimeLimit   Optional[int64]
	IsOnetime   Optional[bool]
	Repeatable  Optional[bool]
	Finished    Optional[bool]
}
