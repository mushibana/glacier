package repository

import (
	"context"

	"github.com/mushibana/glacier/backend/internal/domain"
)

type TaskRepository interface {
	List(ctx context.Context) ([]domain.Task, error)
	Get(ctx context.Context, id int64) (domain.Task, error)
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id int64) error
	Move(ctx context.Context, id int64, beforeID *int64) ([]domain.Task, error)
}
