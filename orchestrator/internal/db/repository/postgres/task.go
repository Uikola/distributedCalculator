package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// TaskRepository структура для работы с репозиторием бд(таблицей tasks).
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository возвращает новую инстанцию TaskRepository.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// AddTask добавляет новую задачу в бд и возвращает ошибку.
func (r TaskRepository) AddTask(ctx context.Context, task entity.Task) (int64, error) {
	query := `
	INSERT INTO tasks(id, expression, status, created_at) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	row := r.db.QueryRowContext(ctx, query, task.ID, task.Expression, task.Status, task.CreatedAt)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Exists проверяет, существует ли задача с указанным id.
func (r TaskRepository) Exists(ctx context.Context, taskID int64) (bool, error) {
	query := `
	SELECT id
	FROM tasks
	WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, taskID)

	err := row.Scan(&taskID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

// ListTask возвращает список задач и ошибку.
func (r TaskRepository) ListTask(ctx context.Context, limit, offset int) ([]entity.Task, error) {
	query := `
	SELECT id, expression, status, created_at, calculated_at
	FROM tasks
	LIMIT $1
	OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	var tasks []entity.Task
	var id int64
	var expression string
	var status entity.ExpressionStatus
	var createdAt time.Time
	var calculatedAt *time.Time

	for rows.Next() {
		err = rows.Scan(&id, &expression, &status, &createdAt, &calculatedAt)
		if err != nil {
			return nil, err
		}

		task := entity.Task{
			ID:           id,
			Expression:   expression,
			Status:       status,
			CreatedAt:    createdAt,
			CalculatedAt: calculatedAt,
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору и ошибку.
func (r TaskRepository) GetTask(ctx context.Context, id int64) (entity.Task, error) {
	query := `
	SELECT id, expression, status, created_at, calculated_at, calculated_by
	FROM tasks
	WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var expression string
	var status entity.ExpressionStatus
	var createdAt time.Time
	var calculatedAt *time.Time
	var calculatedBy *int64

	err := row.Scan(&id, &expression, &status, &createdAt, &calculatedAt, &calculatedBy)
	if err != nil {
		return entity.Task{}, err
	}

	return entity.Task{
		ID:           id,
		Expression:   expression,
		Status:       status,
		CreatedAt:    createdAt,
		CalculatedAt: calculatedAt,
		CalculatedBy: calculatedBy,
	}, nil
}

// UpdateComputingResource обновляет вычислительный ресурс задачи(помечает, что задачу выполняет тот то ресурс)
// и возвращает ошибку.
func (r TaskRepository) UpdateComputingResource(ctx context.Context, taskID, cResourceID int64) error {
	query := `
	UPDATE tasks
	SET calculated_by = $1
	WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, cResourceID, taskID)
	if err != nil {
		return err
	}

	return nil
}

// CompleteTask помечает задачу выполненной и возвращает ошибку.
func (r TaskRepository) CompleteTask(ctx context.Context, taskID int64) error {
	query := `
	UPDATE tasks
	SET status = $1,
	    calculated_at = $2
	WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, entity.OK, time.Now(), taskID)
	if err != nil {
		return err
	}

	return nil
}

// GetTaskByCResource возвращает задачу по её вычислительному ресурсу и возвращает ошибку.
func (r TaskRepository) GetTaskByCResource(ctx context.Context, cResourceID int64) (entity.Task, error) {
	query := `
	SELECT id, expression
	FROM tasks
	WHERE status = $1 AND calculated_by = $2`

	row := r.db.QueryRowContext(ctx, query, entity.InProgress, cResourceID)
	var taskID int64
	var expression string

	err := row.Scan(&taskID, &expression)
	if err != nil {
		return entity.Task{}, err
	}

	return entity.Task{
		ID:         taskID,
		Expression: expression,
	}, nil
}
