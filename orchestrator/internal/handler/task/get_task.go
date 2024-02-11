package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// GetTask ручка, которая возвращает задачу по её идентификатору.
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

	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.Error().Err(err).Msg("can't encode the task")
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
}
