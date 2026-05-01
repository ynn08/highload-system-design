package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/user/highload-system-design/order-service/internal/infrastructure/messaging"
	"github.com/user/highload-system-design/order-service/internal/infrastructure/persistence"
	"github.com/user/highload-system-design/order-service/internal/interface/http"
	"github.com/user/highload-system-design/order-service/internal/usecase"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/order_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
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
	log.Println("Database migrations applied successfully")

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":6379",
	})

	brokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	orderRepo := persistence.NewPostgresRepository(db)
	cartRepo := persistence.NewRedisCartRepository(rdb)

	createOrderUC := usecase.NewCreateOrderUseCase(orderRepo, orderRepo, cartRepo)
	manageCartUC := usecase.NewManageCartUseCase(cartRepo)
	handler := http.NewOrderHandler(createOrderUC, manageCartUC)

	// Start Outbox Processor
	outboxProc := messaging.NewOutboxProcessor(orderRepo, []string{brokers})
	go outboxProc.Start(context.Background())

	// Start Payment Result Listener
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   "payment-processed",
		StartOffset: kafka.FirstOffset,
	})
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Order status consumer error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			var event struct {
				OrderID string `json:"order_id"`
				Status  string `json:"status"`
			}
			if err := json.Unmarshal(m.Value, &event); err == nil {
				newStatus := "PAID"
				if event.Status != "PAID" && event.Status != "SUCCESS" {
					newStatus = "CANCELLED"
				}
				_, _ = db.Exec("UPDATE orders SET status = $1 WHERE order_id = $2", newStatus, event.OrderID)
			}
		}
	}()

	r := gin.Default()
	r.POST("/api/v1/orders", handler.CreateOrder)
	r.POST("/api/v1/carts/:customerId", handler.AddToCart)
	r.GET("/api/v1/carts/:customerId", handler.GetCart)
	r.GET("/health", func(c *gin.Context) { c.Status(200) })

	log.Println("Order service starting on :8080")
	r.Run(":8080")
}
