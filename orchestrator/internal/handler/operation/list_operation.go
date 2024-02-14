package operation

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// ListOperation ручка, которая в качестве ответа выдаёт список доступных операций.
//
// # ListOperation
//
//	@Summary		Выводит список доступных операций.
//	@Description	Выводит список доступных операций и время их работы.
//	@Tags			operations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.Operation
//	@Failure		500	{object}	string
//
//	@Router			/operations [get]
func (h Handler) ListOperation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	operations, err := h.useCase.ListOperation(ctx)
	if err != nil {
		log.Error().Err(err).Msg("can't get a list of operations")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(operations)
}
