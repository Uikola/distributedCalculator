package cresource_usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Uikola/yandexDAEC/orchestrator/internal/entity"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka/consumer"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka/producer"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// AddCResource в отдельной горутине вызывает прослушку heartbeat сообщений(в случае ошибки удаляет
// созданную запись, вызвав метод DeleteCResource) и вызывает метод репозитория AddCResource.
func (uc UseCaseImpl) AddCResource(ctx context.Context, cResource entity.ComputingResource) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		ctx := context.Background()
		err := consumer.StartHeartbeatConsumer(cResource.Name)
		if err != nil {
			createdCResource, err := uc.cResourceRepository.GetCResourceByName(ctx, cResource.Name)
			if err != nil {
				return err
			}

			if createdCResource.Occupied {
				log.Info().Msg("try to resend message")
				task, err := uc.taskRepository.GetTaskByCResource(ctx, createdCResource.ID)
				if err != nil {
					return err
				}

				resources, err := uc.cResourceRepository.ListFreeCResource(ctx)
				switch {
				case errors.Is(err, sql.ErrNoRows):
					time.Sleep(50 * time.Second)
					resources, err = uc.cResourceRepository.ListFreeCResource(ctx)
				case err != nil:
					return err
				}

				resource := resources[rand.Intn(len(resources))]
				err = uc.cResourceRepository.OccupyCResource(ctx, resource.ID, task.Expression)
				if err != nil {
					return err
				}

				err = uc.taskRepository.UpdateComputingResource(ctx, task.ID, resource.ID)
				if err != nil {
					return err
				}

				p, err := producer.NewKafkaProducer(kafka.Brokers)
				if err != nil {
					return err
				}

				msg := fmt.Sprintf("%s:%d", task.Expression, task.ID)
				p.SendMessage(kafka.ExpressionTopic, msg, resource.Name)
			}
			err = uc.cResourceRepository.DeleteCResource(context.Background(), cResource.Name)
			fmt.Println(cResource.Name)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return uc.cResourceRepository.AddCResource(ctx, cResource)
}

func (uc UseCaseImpl) ListCResources(ctx context.Context) (map[string]*string, error) {
	cResources, err := uc.cResourceRepository.ListCResource(ctx)
	if err != nil {
		return nil, err
	}

	pairs := make(map[string]*string)

	for _, cResource := range cResources {
		pairs[cResource.Name] = cResource.Task
	}

	return pairs, nil
}
