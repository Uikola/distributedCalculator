package operation

import (
	"context"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// UseCase интерфейс для взаимодействия с юс кейсом операций.
type UseCase interface {
	ListOperation(ctx context.Context) ([]entity.Operation, error)
	UpdateOperationTime(ctx context.Context, operation entity.Operation) error
}

// Handler структура для работы с operation handler.
type Handler struct {
	useCase UseCase
}

// New возвращает новую инстанцию Handler.
func New(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}
