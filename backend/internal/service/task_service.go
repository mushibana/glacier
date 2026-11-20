package service

import (
	"context"
	"time"

	"github.com/mushibana/glacier/backend/internal/domain"
	"github.com/mushibana/glacier/backend/internal/repository"
)

type TaskService interface {
	List(ctx context.Context) ([]domain.Task, error)
	Get(ctx context.Context, id int64) (domain.Task, error)
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	Update(ctx context.Context, task domain.UpdateTaskParams) (time.Time, error)
	Delete(ctx context.Context, id int64) error
	Move(ctx context.Context, id int64, beforeID *int64) ([]domain.Task, error)
}

type TaskServiceImplementation struct {
	tasks repository.TaskRepository
}

func NewTaskServiceImplementation(taskRepo repository.TaskRepository) *TaskServiceImplementation {
	return &TaskServiceImplementation{tasks: taskRepo}
}

func (s *TaskServiceImplementation) List(ctx context.Context) ([]domain.Task, error) {
	return s.tasks.List(ctx)
}

func (s *TaskServiceImplementation) Get(ctx context.Context, id int64) (domain.Task, error) {
	return s.tasks.Get(ctx, id)
}

func (s *TaskServiceImplementation) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	return s.tasks.Create(ctx, task)
}

func (s *TaskServiceImplementation) Update(ctx context.Context, params domain.UpdateTaskParams) (time.Time, error) {
	return s.tasks.Update(ctx, params)
}

func (s *TaskServiceImplementation) Delete(ctx context.Context, id int64) error {
	return s.tasks.Delete(ctx, id)
}

func (s *TaskServiceImplementation) Move(ctx context.Context, id int64, beforeID *int64) ([]domain.Task, error) {
	return s.tasks.Move(ctx, id, beforeID)
}
