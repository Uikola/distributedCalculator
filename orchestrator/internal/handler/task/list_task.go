package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// ListTask ручка, которая возвращает список задач с лимитом и офсетом.
//
//	@Summary		Получает список задач с лимитом и офсетом.
//	@Description	Получает список задач с лимитом и офсетом.
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	true	"Лимит"
//	@Param			offset	query		int	true	"Смещение"
//	@Success		200		{array}		entity.Task
//	@Failure		500		{object}	string
//
//	@Router			/tasks [get]
func (h Handler) ListTask(w http.ResponseWriter, r *http.Request) {
	var offset, limit int
	ctx := r.Context()

	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 4
	}

	tasks, err := h.useCase.ListTask(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("can't get list of tasks")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
}
