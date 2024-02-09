package postgres

import (
	"context"
	"database/sql"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// OperationRepository структура для работы с репозиторием бд(таблицей operations).
type OperationRepository struct {
	db *sql.DB
}

// NewOperationRepository возвращает новую инстанцию OperationRepository.
func NewOperationRepository(db *sql.DB) *OperationRepository {
	return &OperationRepository{db: db}
}

// ListOperation возвращает список доступных операций с временем их выполнения и ошибку.
func (r OperationRepository) ListOperation(ctx context.Context) ([]entity.Operation, error) {
	query := `
	SELECT name, duration
	FROM operations`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var operations []entity.Operation
	var name string
	var duration int

	for rows.Next() {
		err = rows.Scan(&name, &duration)
		if err != nil {
			return nil, err
		}

		operation := entity.Operation{
			Name:     name,
			Duration: duration,
		}
		operations = append(operations, operation)
	}

	return operations, nil
}

// UpdateOperationTime обновляет время операции с указанным именем и возвращает ошибку.
func (r OperationRepository) UpdateOperationTime(ctx context.Context, operation entity.Operation) error {
	query := `
	UPDATE operations
	SET duration = $1
	WHERE name = $2`

	_, err := r.db.ExecContext(ctx, query, operation.Duration, operation.Name)
	if err != nil {
		return err
	}

	return nil
}
