package repository

import (
	"context"
	"database/sql"
)

// Repository структура для работы с репозиторием бд(таблицей operations).
type Repository struct {
	db *sql.DB
}

// NewRepository возвращает новую инстанцию Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ListOperation возвращает список доступных операций с временем их выполнения.
func (r Repository) ListOperation(ctx context.Context) (map[string]int, error) {
	query := `
	SELECT name, duration
	FROM operations`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	operations := make(map[string]int)
	var name string
	var duration int

	for rows.Next() {
		err = rows.Scan(&name, &duration)
		if err != nil {
			return nil, err
		}

		operations[name] = duration
	}

	return operations, nil
}
