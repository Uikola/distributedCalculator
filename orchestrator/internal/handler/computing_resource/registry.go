package computing_resource

import (
	"encoding/json"
	"fmt"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
	"github.com/rs/zerolog/log"
	"net/http"
)

// RegistryReqBody структура для работы с телом запроса.
type RegistryReqBody struct {
	Name string `json:"name"`
}

// Registry ручка, которая регистрирует новый сервис-вычислитель.
func (h Handler) Registry(w http.ResponseWriter, r *http.Request) {
	var request RegistryReqBody
	ctx := r.Context()
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("can't decode the request")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	cResource := entity.ComputingResource{
		Name: request.Name,
	}
	err := h.useCase.AddCResource(ctx, cResource)
	if err != nil {
		log.Error().Err(err).Msg("can't add the computing resource")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
