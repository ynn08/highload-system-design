package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/user/highload-system-design/payment-service/internal/domain"
	"github.com/user/highload-system-design/payment-service/internal/infrastructure/messaging"
	"github.com/user/highload-system-design/payment-service/internal/infrastructure/persistence"
	"github.com/user/highload-system-design/payment-service/internal/usecase"
)

type OrderCreatedEvent struct {
	OrderID      string  `json:"order_id"`
	RestaurantID string  `json:"restaurant_id"`
	TotalAmount  float64 `json:"total_amount"`
}

// Gateway implementation for processing actual payments
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

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatal("Failed to create postgres driver:", err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", driver)
		if err != nil {
			log.Fatal("Failed to load migrations:", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations:", err)
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
		StartOffset: kafka.FirstOffset,
	})

	log.Printf("Payment service starting, listening on %s for topic order-created", brokers)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var event OrderCreatedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		log.Printf("Processing payment for order %s...", event.OrderID)
		err = processPaymentUC.Execute(context.Background(), event.OrderID, event.RestaurantID, event.TotalAmount)
		if err != nil {
			log.Printf("ProcessPayment error: %v", err)
		} else {
			log.Printf("Successfully processed payment for order %s", event.OrderID)
		}
	}
}
