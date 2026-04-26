package domain

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusSuccess   PaymentStatus = "SUCCESS"
	StatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID            string
	OrderID       string
	Amount        float64
	Status        PaymentStatus
	IdempotencyKey string
	CreatedAt     time.Time
}

type PaymentRepository interface {
	Save(ctx context.Context, payment *Payment) error
	GetByIdempotencyKey(ctx context.Context, key string) (*Payment, error)
}

type MessagePublisher interface {
	PublishPaymentProcessed(ctx context.Context, orderID string, restaurantID string, status string) error
}

type PaymentGateway interface {
	Charge(ctx context.Context, amount float64) (bool, error)
}
