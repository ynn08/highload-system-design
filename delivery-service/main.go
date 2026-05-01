package main

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/user/highload-system-design/delivery-service/internal/domain"
	"github.com/user/highload-system-design/delivery-service/internal/infrastructure/persistence"
	"github.com/user/highload-system-design/delivery-service/internal/interface/http"
	"github.com/user/highload-system-design/delivery-service/internal/usecase"
)

type PaymentProcessedEvent struct {
	OrderID      string `json:"order_id"`
	RestaurantID string `json:"restaurant_id"`
	Status       string `json:"status"`
}

func main() {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	repo := persistence.NewInMemoryRepository()
	assignUC := usecase.NewAssignCourierUseCase(repo, repo)
	saveCourierUC := usecase.NewSaveCourierUseCase(repo)
	handler := http.NewDeliveryHandler(assignUC, saveCourierUC, repo)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   domain.KafkaTopicPaymentProcessed,
		StartOffset: kafka.FirstOffset,
	})

	log.Printf("Delivery service starting, consumer listening on %s for topic %s", brokers, domain.KafkaTopicPaymentProcessed)

	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Consumer error: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var event PaymentProcessedEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Unmarshal error: %v", err)
				continue
			}

			if event.Status == "PAID" {
				dest := domain.Location{Lat: 55.75, Lon: 37.62}
				err := assignUC.Execute(context.Background(), event.OrderID, event.RestaurantID, dest)
				if err != nil {
					log.Printf("Assign error: %v", err)
				}
			}
		}
	}()

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/orders/:orderId/status", handler.GetStatus)
		v1.POST("/couriers", handler.SaveCourier)
	}
	r.GET("/health", func(c *gin.Context) { c.Status(200) })

	log.Println("Delivery service HTTP server starting on :8080")
	r.Run(":8080")
}
