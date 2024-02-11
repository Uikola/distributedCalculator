package consumer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/errorz"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/kafka"
)

// StartResultConsumer запускает kafka consumer, который прослушивает сообщения с результатами вычислений.
func StartResultConsumer(key string) (string, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(kafka.Brokers, config)
	if err != nil {
		return "", err
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(kafka.ResultTopic, 0, sarama.OffsetOldest)
	if err != nil {
		return "", err
	}
	defer partitionConsumer.Close()

	// горутина для обработки полученных result сообщений
	resCh := make(chan string)
	errCh := make(chan error)
	fmt.Println(key)
	go func() {
		timer := time.NewTimer(1 * time.Second)
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				for string(msg.Key) != key {
					select {
					case msg = <-partitionConsumer.Messages():
					case <-timer.C:
						resCh <- ""
						errCh <- errorz.ErrResultNotReady
						return
					}
				}
				resCh <- string(msg.Value)
				errCh <- nil
				return
			case <-timer.C:
				resCh <- ""
				errCh <- errorz.ErrResultNotReady
				return
			}
		}
	}()
	return <-resCh, <-errCh
}
