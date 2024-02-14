package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// GetTask ручка, которая возвращает задачу по её идентификатору.
//
// # GetTask
//
//	@Summary		Получает задачу по её идентификатору.
//	@Description	Получает задачу по её идентификатору.
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task_id	path		int	true	"Идентификатор задачи"
//	@Success		200		{object}	entity.Task
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//
//	@Router			/tasks/{task_id} [get]
func (h Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("invalid task id")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	task, err := h.useCase.GetTask(ctx, int64(taskID))
	if err != nil {
		log.Error().Err(err).Msg("can't get task")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(task)
}
