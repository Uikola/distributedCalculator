package task_usecase

import (
	"context"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// TaskRepository интерфейс репозитория для работы с задачами.
type TaskRepository interface {
	AddTask(ctx context.Context, task entity.Task) (int64, error)
	Exists(ctx context.Context, taskID int64) (bool, error)
	ListTask(ctx context.Context, limit, offset int) ([]entity.Task, error)
	GetTask(ctx context.Context, id int64) (entity.Task, error)
	UpdateComputingResource(ctx context.Context, taskID, cResourceID int64) error
	CompleteTask(ctx context.Context, taskID int64) error
	ErrorTask(ctx context.Context, taskID int64) error
}

// CResourceRepository интерфейс репозитория для работы с вычислительными ресурсами.
type CResourceRepository interface {
	ListFreeCResource(ctx context.Context) ([]entity.ComputingResource, error)
	OccupyCResource(ctx context.Context, cResourceID int64) error
	FreeCResource(ctx context.Context, cResourceID int64) error
	GetCResource(ctx context.Context, id int64) (entity.ComputingResource, error)
}

// UseCaseImpl структура для работы с имплементацией юс кейса задач.
type UseCaseImpl struct {
	taskRepository      TaskRepository
	cResourceRepository CResourceRepository
}

// New возвращает новую инстанцию UseCaseImpl.
func New(taskRepository TaskRepository, cResourceRepository CResourceRepository) *UseCaseImpl {
	return &UseCaseImpl{
		taskRepository:      taskRepository,
		cResourceRepository: cResourceRepository,
	}
}
