package messaging

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/user/highload-system-design/order-service/internal/domain"
)

type OutboxProcessor struct {
	repo   domain.OutboxRepository
	writer *kafka.Writer
}

func NewOutboxProcessor(repo domain.OutboxRepository, brokers []string) *OutboxProcessor {
	return &OutboxProcessor{
		repo: repo,
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "order-created",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, err := p.repo.FindPending(ctx)
			if err != nil {
				log.Printf("Outbox error: %v", err)
				continue
			}

			for _, e := range events {
				err := p.writer.WriteMessages(ctx, kafka.Message{
					Value: []byte(e.Payload),
				})
				if err != nil {
					log.Printf("Kafka publish error: %v", err)
					continue
				}
				_ = p.repo.MarkAsProcessed(ctx, e.ID)
			}
		}
	}
}
