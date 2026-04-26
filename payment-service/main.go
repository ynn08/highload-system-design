package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/user/highload-system-design/payment-service/internal/domain"
	"github.com/user/highload-system-design/payment-service/internal/infrastructure/messaging"
	"github.com/user/highload-system-design/payment-service/internal/infrastructure/persistence"
	"github.com/user/highload-system-design/payment-service/internal/usecase"
	"log"
	"os"
)

type OrderCreatedEvent struct {
	OrderID      string  `json:"orderId"`
	RestaurantID string  `json:"restaurantId"`
	TotalAmount  float64 `json:"totalAmount"`
}

// Mock Gateway
type mockGateway struct{}

func (g *mockGateway) Charge(ctx context.Context, amount float64) (bool, error) {
	return true, nil
}

func main() {
	dbURL := os.Getenv("DB_URL")
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	var repo domain.PaymentRepository
	if dbURL != "" {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		repo = persistence.NewPostgresPaymentRepository(db)
	} else {
		log.Println("DB_URL not set, using mock repo")
	}

	publisher := messaging.NewKafkaPublisher([]string{brokers}, "payment-processed")
	gateway := &mockGateway{}
	processPaymentUC := usecase.NewProcessPaymentUseCase(repo, publisher, gateway)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   "order-created",
		GroupID: "payment-group",
	})

	log.Println("Payment service started...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		var event OrderCreatedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Error unmarshaling: %v", err)
			continue
		}

		err = processPaymentUC.Execute(context.Background(), event.OrderID, event.RestaurantID, event.TotalAmount)
		if err != nil {
			log.Printf("Error processing payment: %v", err)
		}
	}
}
