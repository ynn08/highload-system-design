package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/user/highload-system-design/payment-service/internal/domain"
	"log"
	"time"
)

type ProcessPaymentUseCase struct {
	repo      domain.PaymentRepository
	publisher domain.MessagePublisher
	gateway   domain.PaymentGateway
}

func NewProcessPaymentUseCase(repo domain.PaymentRepository, publisher domain.MessagePublisher, gateway domain.PaymentGateway) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{repo: repo, publisher: publisher, gateway: gateway}
}

func (u *ProcessPaymentUseCase) Execute(ctx context.Context, orderID string, restaurantID string, amount float64) error {
	// Non-trivial logic: Idempotency check
	// Using orderID as idempotency key for this demo
	existing, err := u.repo.GetByIdempotencyKey(ctx, orderID)
	if err == nil && existing != nil {
		log.Printf("Payment for order %s already processed with status %s", orderID, existing.Status)
		return u.publisher.PublishPaymentProcessed(ctx, orderID, restaurantID, string(existing.Status))
	}

	payment := &domain.Payment{
		ID:             uuid.New().String(),
		OrderID:        orderID,
		Amount:         amount,
		Status:         domain.StatusPending,
		IdempotencyKey: orderID,
		CreatedAt:      time.Now(),
	}

	success, err := u.gateway.Charge(ctx, amount)
	if err != nil || !success {
		payment.Status = domain.StatusFailed
	} else {
		payment.Status = domain.StatusSuccess
	}

	if err := u.repo.Save(ctx, payment); err != nil {
		return err
	}

	statusStr := "PAID"
	if payment.Status == domain.StatusFailed {
		statusStr = "FAILED"
	}

	return u.publisher.PublishPaymentProcessed(ctx, orderID, restaurantID, statusStr)
}
