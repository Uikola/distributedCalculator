package producer

import (
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

// KafkaProducer структура для работы с продюсером Kafka
type KafkaProducer struct {
	Producer sarama.AsyncProducer
}

// NewKafkaProducer возвращает инстанцию асинхронного KafkaProducer
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Error().Err(err).Msg("Failed to send message to Kafka")
		}
	}()

	return &KafkaProducer{Producer: producer}, nil
}

// SendMessage отправляет сообщение в Kafka
func (kp *KafkaProducer) SendMessage(topic, message, key string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	kp.Producer.Input() <- msg
}
