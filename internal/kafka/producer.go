package kafka

import (
	"context"
	"encoding/json"
	"log"
	"github.com/gulovv/neuron_bot/models"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafkaWriter() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "messages",
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("[Kafka] Writer инициализирован")
}

func SendToKafka(ctx context.Context, msg models.Message) error {
	if kafkaWriter == nil {
		InitKafkaWriter()
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Username),
		Value: data,
	})
	return err
}