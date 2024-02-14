package computing_resource

import (
	"context"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// UseCase интерфейс для работы с юс кейсом вычислительных ресурсов.
type UseCase interface {
	AddCResource(ctx context.Context, cResource entity.ComputingResource) error
	ListCResources(ctx context.Context) (map[string]*string, error)
}

// Handler структура для работы с хендлером вычислительных ресурсов.
type Handler struct {
	useCase UseCase
}

// New возвращает новую инстанцию Handler.
func New(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}
