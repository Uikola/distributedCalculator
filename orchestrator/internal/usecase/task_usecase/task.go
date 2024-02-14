package task_usecase

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/errorz"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka/consumer"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka/producer"
)

// AddTask отправляет арифметическое выражение в топик "expressions" и вызывает метод репозитория AddTask.
func (uc UseCaseImpl) AddTask(ctx context.Context, task entity.Task) (int64, error) {
	exists, err := uc.taskRepository.Exists(ctx, task.ID)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, errorz.ErrTaskAlreadyExists
	}

	id, err := uc.taskRepository.AddTask(ctx, task)
	if err != nil {
		return 0, err
	}

	resources, err := uc.cResourceRepository.ListFreeCResource(ctx)
	if err != nil {
		err = uc.taskRepository.ErrorTask(ctx, id)
		return 0, err
	}
	if len(resources) == 0 {
		err = uc.taskRepository.ErrorTask(ctx, id)
		return 0, errorz.ErrNoAvailableResources
	}

	resource := resources[rand.Intn(len(resources))]
	err = uc.cResourceRepository.OccupyCResource(ctx, resource.ID, task.Expression)
	if err != nil {
		err = uc.taskRepository.ErrorTask(ctx, id)
		return 0, err
	}

	err = uc.taskRepository.UpdateComputingResource(ctx, id, resource.ID)
	if err != nil {
		err = uc.taskRepository.ErrorTask(ctx, id)
		return 0, err
	}

	p, err := producer.NewKafkaProducer(kafka.Brokers)
	if err != nil {
		err = uc.taskRepository.ErrorTask(ctx, id)
		return 0, err
	}
	msg := fmt.Sprintf("%s:%d", task.Expression, id)
	p.SendMessage(kafka.ExpressionTopic, msg, resource.Name)
	return id, nil
}

// ListTask вызывает метод репозитория ListTask.
func (uc UseCaseImpl) ListTask(ctx context.Context, limit, offset int) ([]entity.Task, error) {
	return uc.taskRepository.ListTask(ctx, limit, offset)
}

// GetTask вызывает метод репозитория GetTask.
func (uc UseCaseImpl) GetTask(ctx context.Context, id int64) (entity.Task, error) {
	return uc.taskRepository.GetTask(ctx, id)
}

// GetResult получает результат из очереди и помечает вычислительных ресурс свободным.
func (uc UseCaseImpl) GetResult(ctx context.Context, id int64) (string, error) {
	task, err := uc.taskRepository.GetTask(ctx, id)
	if err != nil {
		return "", err
	}

	cResource, err := uc.cResourceRepository.GetCResource(ctx, *task.CalculatedBy)
	if err != nil {
		return "", err
	}
	res, err := consumer.StartResultConsumer(strconv.FormatInt(id, 10))
	if err != nil {
		return "", err
	}
	err = uc.taskRepository.CompleteTask(ctx, id)
	if err != nil {
		return "", err
	}

	err = uc.cResourceRepository.FreeCResource(ctx, cResource.ID)
	if err != nil {
		return "", nil
	}

	return res, nil
}
