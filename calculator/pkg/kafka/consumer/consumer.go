package consumer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/Uikola/yandexDAEC/calculator/pkg/kafka"
	"github.com/Uikola/yandexDAEC/calculator/pkg/polish_notation"
	"github.com/rs/zerolog/log"
)

// StartConsumer запускает kafka producer, который пишет heartbeat сообщения в топик "heartbeat", и kafka
// consumer, который обрабатывает сообщения из топика "expressions" и записывает их в топик "results".
func StartConsumer(name string, operations map[string]int) {
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(kafka.Brokers, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka consumer")
	}
	defer consumer.Close()

	producer, err := sarama.NewAsyncProducer(kafka.Brokers, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka async producer")
	}
	defer producer.Close()

	partitionConsumer, err := consumer.ConsumePartition(kafka.ExpressionTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create partition consumer")
	}
	defer partitionConsumer.Close()

	go sendHeartbeats(name, producer)

	for msg := range partitionConsumer.Messages() {
		if string(msg.Key) == name {

			log.Info().Msg(fmt.Sprintf("got message from %s, expression: %s", string(msg.Key), string(msg.Value)))
			l := strings.Split(string(msg.Value), ":")
			expression, taskID := l[0], l[1]

			// вычисление выражения
			log.Info().Msg("calculating result")
			rpn := polish_notation.ConvertToRPN(expression)
			result := polish_notation.EvalRPN(rpn, operations)
			log.Info().Msg(fmt.Sprintf("calculation done: %d", result))

			msg := &sarama.ProducerMessage{
				Topic: kafka.ResultTopic,
				Key:   sarama.StringEncoder(taskID),
				Value: sarama.StringEncoder(strconv.Itoa(result)),
			}

			producer.Input() <- msg
		}
	}
}

// sendHeartbeats создаёт kafka producer, который шлёт heartbeat сообщения.
func sendHeartbeats(name string, producer sarama.AsyncProducer) {
	ticker := time.NewTicker(5 * time.Second) // Отправка heartbeat каждые 5 секунд
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Отправка "heartbeat" сообщения
			msg := &sarama.ProducerMessage{
				Topic: kafka.HeartbeatTopic,
				Key:   sarama.StringEncoder(name),
				Value: sarama.StringEncoder("heartbeat"),
			}
			producer.Input() <- msg
		}
	}
}
