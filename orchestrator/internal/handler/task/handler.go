package task

import (
	"context"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// UseCase интерфейс для взаимодействия с юс кейсом задач.
type UseCase interface {
	AddTask(ctx context.Context, task entity.Task) (int64, error)
	ListTask(ctx context.Context, limit, offset int) ([]entity.Task, error)
	GetTask(ctx context.Context, id int64) (entity.Task, error)
	GetResult(ctx context.Context, id int64) (string, error)
}

// Handler структура для работы с task handler.
type Handler struct {
	useCase UseCase
}

// New возвращает инстанцию Handler.
func New(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}
