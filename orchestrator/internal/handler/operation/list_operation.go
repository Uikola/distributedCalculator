package operation

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// ListOperation ручка, которая в качестве ответа выдаёт список доступных операций.
func (h Handler) ListOperation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	operations, err := h.useCase.ListOperation(ctx)
	if err != nil {
		log.Error().Err(err).Msg("can't get a list of operations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(operations); err != nil {
		log.Error().Err(err).Msg("can't encode a list of operations")
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
}
