package operation_usecase

import (
	"context"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// OperationRepository интерфейс репозитория для работы с арифметическими операциями.
type OperationRepository interface {
	ListOperation(ctx context.Context) ([]entity.Operation, error)
	UpdateOperationTime(ctx context.Context, operation entity.Operation) error
}

// UseCaseImpl структура для работы с имплементацией юс кейса арифметических операций.
type UseCaseImpl struct {
	repository OperationRepository
}

// New возвращает новую инстанцию UseCaseImpl.
func New(repo OperationRepository) *UseCaseImpl {
	return &UseCaseImpl{repository: repo}
}
