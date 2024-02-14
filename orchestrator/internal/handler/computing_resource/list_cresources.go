package computing_resource

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ListCResources выводит список вычислительных ресурсов с задачами, которые на них выполняются.
//
// # LoadCResources
//
//	@Summary		Выводит список вычислительных мощностей с задачами, которые на них выполняются.
//	@Description	Выводит список вычислительных мощностей с задачами, которые на них выполняются.
//	@Tags			computing_resources
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.ComputingResource
//	@Failure		500	{object}	string
//
//	@Router			/c_resources [get]
func (h Handler) ListCResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cResources, err := h.useCase.ListCResources(ctx)
	if err != nil {
		log.Error().Err(err).Msg("can't get list of computing resources")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(cResources)
}
