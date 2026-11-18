package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mushibana/glacier/backend/internal/domain"
)

type TaskRepositoryImplementation struct {
	db *sql.DB
}

func NewTaskRepositoryImplementation(db *sql.DB) TaskRepositoryImplementation {
	return TaskRepositoryImplementation{db: db}
}

func (r *TaskRepositoryImplementation) List(ctx context.Context) ([]domain.Task, error) {
	const query = `
select 
	id,
	parent_id,
	order_num,
	name,
	stream_id,
	is_onetime,
	repeatable,
	finished
from tasks
where deleted_at is null
order by order_num`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query task list: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var i tasksRow

		if err = rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.OrderNum,
			&i.Name,
			&i.StreamID,
			&i.IsOnetime,
			&i.Repeatable,
			&i.Finished,
		); err != nil {
			return nil, fmt.Errorf("scan to tasksRow: %w", err)
		}
		task, err := toDomain(i)
		if err != nil {
			return nil, fmt.Errorf("transform to task domain: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close rows: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepositoryImplementation) Get(ctx context.Context, id int64) (domain.Task, error) {
	const query = `
select 
	id,
	parent_id,
	order_num,
	name,
	description,
	stream_id,
	start_at,
	deadline_at,
	timelimit,
	is_onetime,
	repeatable,
	finished,
	created_at,
	updated_at
from tasks
where deleted_at is null and id = ?`

	var i tasksRow
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID,
		&i.ParentID,
		&i.OrderNum,
		&i.Name,
		&i.Description,
		&i.StreamID,
		&i.StartAt,
		&i.DeadlineAt,
		&i.TimeLimit,
		&i.IsOnetime,
		&i.Repeatable,
		&i.Finished,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("scan to tasksRow: %w", err)
	}

	task, err := toDomain(i)
	if err != nil {
		return domain.Task{}, fmt.Errorf("transform to task domain: %w", err)
	}
	return task, nil
}

func (r *TaskRepositoryImplementation) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	const query = `
insert into tasks(
	parent_id,
	order_num,
	name,
	description,
  	stream_id,
	start_at,
	deadline_at,
  	timelimit,
  	is_onetime,
  	repeatable,
  	finished,
	created_at
) values (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) returning id`

	t, err := fromDomain(task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("convert from task domain: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		t.ParentID,
		t.OrderNum,
		t.Name,
		t.Description,
		t.StreamID,
		t.StartAt,
		t.DeadlineAt,
		t.TimeLimit,
		t.IsOnetime,
		t.Repeatable,
		t.Finished,
		time.Now().Format(time.RFC3339),
	).Scan(&task.ID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("query insert task, scan new id: %w", err)
	}

	return task, nil
}

func (r *TaskRepositoryImplementation) Update(ctx context.Context, task domain.UpdateTaskParams) (domain.Task, error) {
	if task.ID <= 0 {
		return domain.Task{}, domain.ErrNotFound //TODO: change to invalidParam
	}

	// TODO:
	return domain.Task{}, nil
}

func (r *TaskRepositoryImplementation) Delete(ctx context.Context, id int64) error {
	const query = `update tasks set deleted_at = date('now', 'localtime') where id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query delete task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// func (r *TaskRepositoryImplementation) Move(ctx context.Context, id int64, beforeID *int64) ([]domain.Task, error) {
// }
