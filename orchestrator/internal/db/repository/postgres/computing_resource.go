package postgres

import (
	"context"
	"database/sql"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// CResourceRepository структура для работы с репозиторием бд(таблицей computing_resources).
type CResourceRepository struct {
	db *sql.DB
}

// NewCResourceRepository возвращает новую инстанцию CResourceRepository.
func NewCResourceRepository(db *sql.DB) *CResourceRepository {
	return &CResourceRepository{db: db}
}

// AddCResource добавляет новый вычислительный ресурс в бд и возвращает ошибку.
func (r CResourceRepository) AddCResource(ctx context.Context, cResource entity.ComputingResource) error {
	query := `
	INSERT INTO computing_resources(name)
	VALUES ($1)`

	_, err := r.db.ExecContext(ctx, query, cResource.Name)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCResource удаляет вычислительный ресурс по его имени и возвращает ошибку.
func (r CResourceRepository) DeleteCResource(ctx context.Context, cResourceName string) error {
	query := `
	DELETE FROM computing_resources
	WHERE name = $1`

	_, err := r.db.ExecContext(ctx, query, cResourceName)
	if err != nil {
		return err
	}

	return nil
}

// ListFreeCResource возвращает список свободных вычислительных ресурсов и ошибку.
func (r CResourceRepository) ListFreeCResource(ctx context.Context) ([]entity.ComputingResource, error) {
	query := `
	SELECT * FROM computing_resources
	WHERE occupied = false`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var cResources []entity.ComputingResource
	var id int64
	var name string
	var occupied bool

	for rows.Next() {
		err = rows.Scan(&id, &name, &occupied)
		if err != nil {
			return nil, err
		}
		cResource := entity.ComputingResource{
			ID:       id,
			Name:     name,
			Occupied: occupied,
		}
		cResources = append(cResources, cResource)
	}

	return cResources, nil
}

// OccupyCResource занимает вычислительный ресурс и возвращает ошибку.
func (r CResourceRepository) OccupyCResource(ctx context.Context, cResourceID int64) error {
	query := `
	UPDATE computing_resources
	SET occupied = true
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, cResourceID)
	if err != nil {
		return err
	}

	return nil
}

// FreeCResource освобождает вычислительный ресурс и возвращает ошибку.
func (r CResourceRepository) FreeCResource(ctx context.Context, cResourceID int64) error {
	query := `
	UPDATE computing_resources
	SET occupied = false
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, cResourceID)
	if err != nil {
		return err
	}

	return nil
}

// GetCResource возвращает вычислительный ресурс по его идентификатору и ошибку.
func (r CResourceRepository) GetCResource(ctx context.Context, id int64) (entity.ComputingResource, error) {
	query := `
	SELECT id, name, occupied
	FROM computing_resources
	WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	var name string
	var occupied bool

	err := row.Scan(&id, &name, &occupied)
	if err != nil {
		return entity.ComputingResource{}, nil
	}

	return entity.ComputingResource{
		ID:       id,
		Name:     name,
		Occupied: occupied,
	}, nil
}

// GetCResourceByName возвращает вычислительный ресурс по его имени и ошибку.
func (r CResourceRepository) GetCResourceByName(ctx context.Context, name string) (entity.ComputingResource, error) {
	query := `
	SELECT id, name, occupied
	FROM computing_resources
	WHERE name = $1`

	row := r.db.QueryRowContext(ctx, query, name)
	var id int64
	var occupied bool

	err := row.Scan(&id, &name, &occupied)
	if err != nil {
		return entity.ComputingResource{}, err
	}

	return entity.ComputingResource{
		ID:       id,
		Name:     name,
		Occupied: occupied,
	}, nil
}

// CleanUpCResources освобождает таблицу вычислительных ресурсов.
func (r CResourceRepository) CleanUpCResources(ctx context.Context) error {
	query := `
	DELETE FROM computing_resources`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
