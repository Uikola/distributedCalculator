package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/errorz"
	"github.com/rs/zerolog/log"
)

// AddTaskReqBody структура для работы с телом запроса.
type AddTaskReqBody struct {
	ID         int64  `json:"id"`
	Expression string `json:"expression"`
}

// AddTask ручка, которая добавляет новую задачу.
//
// # AddTask
//
//	@Summary		Добавляет новую задачу и отдаёт её на обработку.
//	@Description	Добавляет новую задачу и отдаёт её на обработку.
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		AddTaskReqBody	true	"Задача"
//	@Success		200		{object}	string
//	@Success		201		{object}	string
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//
//	@Router			/calculate [post]
func (h Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var request AddTaskReqBody
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("can't decode the request")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := ValidateExpression(request.Expression)
	if err != nil {
		log.Error().Err(err).Msg("can't validate the expression")
		http.Error(w, "invalid expression", http.StatusBadRequest)
		return
	}

	task := entity.Task{
		ID:         request.ID,
		Expression: request.Expression,
		Status:     entity.InProgress,
		CreatedAt:  time.Now(),
	}

	id, err := h.useCase.AddTask(ctx, task)
	switch {
	case errors.Is(err, errorz.ErrNoAvailableResources):
		log.Info().Err(err).Msg("there is no available resources")
		http.Error(w, "no available resources", http.StatusBadRequest)
		return
	case errors.Is(err, errorz.ErrTaskAlreadyExists):
		log.Info().Err(err).Msg("task already exists")
		_ = json.NewEncoder(w).Encode("task with this id already exists, try another one")
		return
	case err != nil:
		log.Error().Err(err).Msg("can't add task")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	answer := fmt.Sprintf("you can get the result after a while by sending a request with the following id: %s", int(id))
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(answer)
}

// ValidateExpression проверяет на корректность входящее выражение.
func ValidateExpression(expr string) error {
	re := regexp.MustCompile(`[^0-9+\-*/() ]`)

	if re.MatchString(expr) {
		return errorz.ErrInvalidExpression
	}

	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return errorz.ErrInvalidExpression
	}

	_, err = expression.Evaluate(nil)
	if err != nil {
		return errorz.ErrInvalidExpression
	}

	return nil
}
