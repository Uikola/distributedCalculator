package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Uikola/yandexDAEC/orchestrator/pkg/errorz"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetResult ручка, которая возвращает результат вычисленного выражения.
//
// # GetResult
//
//	@Summary		Получает результат по идентификатору задачи.
//	@Description	Получает результат по идентификатору задачи.
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task_id	path		int	true	"Идентификатор задачи"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//
//	@Router			/results/{task_id} [get]
func (h Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("invalid task id")
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	result, err := h.useCase.GetResult(ctx, int64(taskID))
	switch {
	case errors.Is(err, errorz.ErrResultNotReady):
		log.Info().Msg("the expression will be calculated soon")
		_ = json.NewEncoder(w).Encode(map[string]string{"result": "the expression will be calculated soon"})
		return
	case err != nil:
		log.Error().Err(err).Msg("can't get a result")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"result": result})
}
