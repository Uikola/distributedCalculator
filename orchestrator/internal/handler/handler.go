package handler

import (
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/computing_resource"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/operation"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/task"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Uikola/yandexDAEC/orchestrator/docs"
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
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.Post("/calculate", handler.Task.AddTask)
	router.Get("/tasks/{id}", handler.Task.GetTask)
	router.Get("/tasks", handler.Task.ListTask)
	router.Get("/operations", handler.Operation.ListOperation)
	router.Put("/operations", handler.Operation.UpdateOperationTime)
	router.Post("/registry", handler.cResource.Registry)
	router.Get("/results/{id}", handler.Task.GetResult)
	router.Get("/c_resources", handler.cResource.ListCResources)
}
