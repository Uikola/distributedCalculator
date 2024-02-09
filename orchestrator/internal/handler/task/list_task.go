package task

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// ListTask ручка, которая возвращает список задач с лимитом и офсетом.
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

	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		log.Error().Err(err).Msg("can't encode the list of tasks")
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
}
