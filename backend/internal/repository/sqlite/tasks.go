package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
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
	name,
	order_num
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
			&i.Name,
			&i.OrderNum,
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
		&i.Timelimit,
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

	t := fromDomain(task)

	err := r.db.QueryRowContext(ctx, query,
		t.ParentID,
		t.OrderNum,
		t.Name,
		t.Description,
		t.StreamID,
		t.StartAt,
		t.DeadlineAt,
		t.Timelimit,
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

func (r *TaskRepositoryImplementation) Update(ctx context.Context, params domain.UpdateTaskParams) (time.Time, error) {
	if params.ID <= 0 {
		return time.Time{}, domain.ErrInvalidParam("id")
	}

	taskRow := fromUpdateTaskParams(params)
	var updateFieldsStr string
	args := make([]any, 0)

	switch {
	case params.ParentID.Set:
		updateFieldsStr += "parent_id = ?, "
		args = append(args, taskRow.ParentID)
	case params.Name.Set:
		if params.Name.Value == "" {
			return time.Time{}, domain.ErrInvalidParam("name")
		}
		updateFieldsStr += "name = ?, "
		args = append(args, taskRow.Name)
	case params.Description.Set:
		updateFieldsStr += "description = ?, "
		args = append(args, taskRow.Description)
	case params.StreamID.Set:
		updateFieldsStr += "stream_id = ?, "
		args = append(args, taskRow.StreamID)
	case params.StartAt.Set:
		updateFieldsStr += "start_at = ?, "
		args = append(args, taskRow.StartAt)
	case params.DeadlineAt.Set:
		updateFieldsStr += "deadline_at = ?, "
		args = append(args, taskRow.DeadlineAt)
	case params.Timelimit.Set:
		updateFieldsStr += "timelimit = ?, "
		args = append(args, taskRow.Timelimit)
	case params.IsOnetime.Set:
		updateFieldsStr += "is_onetime = ?, "
		args = append(args, taskRow.IsOnetime)
	case params.Repeatable.Set:
		updateFieldsStr += "repeatable = ?, "
		args = append(args, taskRow.Repeatable)
	case params.Finished.Set:
		updateFieldsStr += "finished = ?, "
		args = append(args, taskRow.Finished)
	}

	if updateFieldsStr == "" || len(args) == 0 {
		return time.Time{}, domain.ErrNoUpdateParams
	}
	updatedAt := time.Now()
	args = append(args, updatedAt.Format(time.RFC3339), params.ID)
	query := fmt.Sprintf("update tasks set %s updated_at = ? where id = ?", updateFieldsStr)

	result, err := r.db.ExecContext(ctx, query, args)
	if err != nil {
		return time.Time{}, fmt.Errorf("query update task, scan: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return time.Time{}, fmt.Errorf("get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return time.Time{}, domain.ErrNotFound
	}

	return updatedAt, nil
}

func (r *TaskRepositoryImplementation) Delete(ctx context.Context, id int64) error {
	const query = `update tasks set deleted_at = ? where id = ?`

	result, err := r.db.ExecContext(ctx, query, time.Now().Format(time.RFC3339), id)
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

func (r *TaskRepositoryImplementation) Move(ctx context.Context, id int64, beforeID *int64) ([]domain.Task, error) {
	tasks, err := r.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get list of tasks: %w", err)
	}

	sortTasksFn := func(a, b domain.Task) int {
		if a.OrderNum < b.OrderNum {
			return -1
		} else if a.OrderNum == b.OrderNum {
			return 0
		} else {
			return 1
		}
	}

	slices.SortFunc(tasks, sortTasksFn)

	// before always should be initialized as -1, so "to = before + 1" in any case either beforeID is nil or not
	from, before := -1, -1
	for i, task := range tasks {
		if beforeID != nil && task.ID == *beforeID {
			before = i
		}
		if task.ID == id {
			from = i
		}

		if from != -1 && (beforeID == nil || before != -1) {
			break
		}
	}

	// moved to end of list
	if before == len(tasks)-1 {
		// we didn't reorder initial task list so "from" is used
		tasks[from].OrderNum = tasks[before].OrderNum + domain.OrderNumDefaultSpace
		return []domain.Task{tasks[from]}, nil
	}

	to := before + 1

	moved := tasks[from]
	if from < to {
		// move task from left to right => shifting list to the left
		copy(tasks[from:to], tasks[from+1:to+1])
	} else {
		// move task from right to left => shifting list to the right
		copy(tasks[to+1:from+1], tasks[to:from])
	}
	tasks[to] = moved

	prevOrderNum := int64(0)
	nextOrderNum := tasks[len(tasks)-1].OrderNum + domain.OrderNumDefaultSpace
	if to-1 >= 0 {
		prevOrderNum = tasks[to-1].OrderNum
	}
	if to+1 <= len(tasks)-1 {
		nextOrderNum = tasks[to+1].OrderNum
	}

	tasks[to].OrderNum = prevOrderNum + (nextOrderNum-prevOrderNum)/2
	if prevOrderNum < tasks[to].OrderNum && nextOrderNum < tasks[to].OrderNum {
		return []domain.Task{tasks[to]}, nil
	}

	// branch when integer nums between neighbours are exhausted (new order number is already in use)
	// so rebalance is required
	for i := range tasks {
		tasks[i].OrderNum = int64(i) * domain.OrderNumDefaultSpace
	}

	return tasks, nil
}
