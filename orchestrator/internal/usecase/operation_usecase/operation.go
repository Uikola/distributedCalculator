package operation_usecase

import (
	"context"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
)

// ListOperation вызывает метод репозитория ListOperation.
func (uc UseCaseImpl) ListOperation(ctx context.Context) ([]entity.Operation, error) {
	return uc.repository.ListOperation(ctx)
}

// UpdateOperationTime вызывает метод репозитория UpdateOperationTime.
func (uc UseCaseImpl) UpdateOperationTime(ctx context.Context, operation entity.Operation) error {
	return uc.repository.UpdateOperationTime(ctx, operation)
}
