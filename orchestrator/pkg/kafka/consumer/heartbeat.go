package consumer

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka"
	"github.com/rs/zerolog/log"
	"time"
)

// StartHeartbeatConsumer запускает kafka consumer, который прослушивает heartbeat сообщения.
func StartHeartbeatConsumer(key string) error {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(kafka.Brokers, config)
	if err != nil {
		return err
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(kafka.HeartbeatTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	// горутина для обработки полученных heartbeat сообщений
	errCh := make(chan error)
	go func() {
		for {
			timer := time.NewTimer(30 * time.Second)
			select {
			case msg := <-partitionConsumer.Messages():
				for string(msg.Key) != key {
					select {
					case msg = <-partitionConsumer.Messages():
					case <-timer.C:
						errCh <- errors.New("computing resource is dead")
						return
					}
				}
				log.Info().Msg(fmt.Sprintf("got heartbeat message from %s", string(msg.Key)))
			case <-timer.C:
				errCh <- errors.New("computing resource is dead")
				return
			}
		}
	}()
	return <-errCh
}
