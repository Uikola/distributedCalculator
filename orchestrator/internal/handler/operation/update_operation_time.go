package operation

import (
	"encoding/json"
	"net/http"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
	"github.com/rs/zerolog/log"
)

// UpdateOperationTimeReqBody структура для работы с телом запроса.
type UpdateOperationTimeReqBody struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

// UpdateOperationTime ручка, которая обновляет время выполнения указанной операции.
func (h Handler) UpdateOperationTime(w http.ResponseWriter, r *http.Request) {
	var request UpdateOperationTimeReqBody
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("can't decode the request")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	operation := entity.Operation{
		Name:     request.Name,
		Duration: request.Duration,
	}

	err := h.useCase.UpdateOperationTime(ctx, operation)
	if err != nil {
		log.Error().Err(err).Msg("can't update operation time")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
