package handler

import (
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/computing_resource"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/operation"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/task"
	"github.com/go-chi/chi/v5"
)

// Handler структура для взаимодействия с хендлером, объединяющем все хендлеры.
type Handler struct {
	Operation *operation.Handler
	Task      *task.Handler
	cResource *computing_resource.Handler
}

// New возвращает инстанцию Handler.
func New(operationHandler *operation.Handler, taskHandler *task.Handler, cResource *computing_resource.Handler) *Handler {
	return &Handler{Operation: operationHandler, Task: taskHandler, cResource: cResource}
}

// Router создаёт всю инфраструктуру и задаёт пути для ручек.
func Router(handler *Handler, router chi.Router) {
	router.Post("/calculate", handler.Task.AddTask)
	router.Get("/task/{id}", handler.Task.GetTask)
	router.Get("/task", handler.Task.ListTask)
	router.Get("/operation", handler.Operation.ListOperation)
	router.Put("/operation", handler.Operation.UpdateOperationTime)
	router.Post("/registry", handler.cResource.Registry)
	router.Get("/result/{id}", handler.Task.GetResult)
}
