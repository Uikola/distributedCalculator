package cresource_usecase

import (
	"context"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// CResourceRepository интерфейс репозитория для работы с вычислительными ресурсами.
type CResourceRepository interface {
	AddCResource(ctx context.Context, cResource entity.ComputingResource) error
	DeleteCResource(ctx context.Context, cResourceName string) error
	ListFreeCResource(ctx context.Context) ([]entity.ComputingResource, error)
	OccupyCResource(ctx context.Context, cResourceID int64, task string) error
	GetCResourceByName(ctx context.Context, name string) (entity.ComputingResource, error)
	ListCResource(ctx context.Context) ([]entity.ComputingResource, error)
}

// TaskRepository интерфейс репозитория для работы с задачами.
type TaskRepository interface {
	AddTask(ctx context.Context, task entity.Task) (int64, error)
	ListTask(ctx context.Context, limit, offset int) ([]entity.Task, error)
	GetTask(ctx context.Context, id int64) (entity.Task, error)
	UpdateComputingResource(ctx context.Context, taskID, cResourceID int64) error
	CompleteTask(ctx context.Context, taskID int64) error
	GetTaskByCResource(ctx context.Context, cResourceID int64) (entity.Task, error)
}

// UseCaseImpl структура для работы с имплементацией юс кейса вычислительных ресурсов.
type UseCaseImpl struct {
	cResourceRepository CResourceRepository
	taskRepository      TaskRepository
}

// New возвращает новую инстанцию UseCaseImpl.
func New(cResourceRepository CResourceRepository, taskRepository TaskRepository) *UseCaseImpl {
	return &UseCaseImpl{cResourceRepository: cResourceRepository, taskRepository: taskRepository}
}
