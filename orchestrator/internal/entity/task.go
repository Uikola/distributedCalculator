package entity

import "time"

// ExpressionStatus описывает статусы выполнения выражения.
type ExpressionStatus string

const (
	Error      ExpressionStatus = "error"
	InProgress ExpressionStatus = "in_progress"
	OK         ExpressionStatus = "ok"
)

// Task структура для работы с задачами. Содержит в себе выражение для вычисления, статус вычисления,
// дату создания и дату завершения вычислений.
type Task struct {
	ID           int64
	Expression   string
	Status       ExpressionStatus
	CreatedAt    time.Time
	CalculatedAt *time.Time
	CalculatedBy *int64
}
