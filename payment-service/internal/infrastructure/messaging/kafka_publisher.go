package messaging

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) PublishPaymentProcessed(ctx context.Context, orderID string, restaurantID string, status string) error {
	msg := map[string]string{
		"orderId":      orderID,
		"restaurantId": restaurantID,
		"status":       status,
	}
	bytes, _ := json.Marshal(msg)
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
}
