package usecase

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/user/highload-system-design/order-service/internal/domain"
)

type CreateOrderUseCase struct {
	orderRepo  domain.OrderRepository
	outboxRepo domain.OutboxRepository
	cartRepo   domain.CartRepository
}

func NewCreateOrderUseCase(or domain.OrderRepository, out domain.OutboxRepository, cr domain.CartRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{orderRepo: or, outboxRepo: out, cartRepo: cr}
}

func (u *CreateOrderUseCase) Execute(ctx context.Context, input domain.Order) (*domain.Order, error) {
	input.OrderID = uuid.New().String()
	input.Status = domain.StatusPending
	input.CreatedAt = time.Now()

	// Dynamic Pricing & ETA
	u.calculateSurgeAndFees(&input)
	input.EstimatedDelivery = time.Now().Add(25 * time.Minute)

	var total float64
	for _, item := range input.Items {
		total += item.Price * float64(item.Quantity)
	}
	input.TotalAmount = total + input.DeliveryFee

	if err := u.orderRepo.Save(ctx, &input); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(input)
	outboxEvent := &domain.OutboxEvent{
		ID:            uuid.New().String(),
		AggregateID:   input.OrderID,
		AggregateType: "ORDER",
		EventType:     "ORDER_CREATED",
		Payload:       string(payload),
		Status:        "PENDING",
		CreatedAt:     time.Now(),
	}

	if err := u.outboxRepo.SaveEvent(ctx, outboxEvent); err != nil {
		return nil, err
	}

	_ = u.cartRepo.DeleteByCustomerID(ctx, input.CustomerID)

	return &input, nil
}

func (u *CreateOrderUseCase) calculateSurgeAndFees(order *domain.Order) {
	hour := time.Now().Hour()
	multiplier := 1.0
	if (hour >= 12 && hour <= 14) || (hour >= 18 && hour <= 21) {
		multiplier = 1.5
	}

	baseFee := 5.0
	distanceKm := 5.0
	perKmRate := 2.0
	fee := (baseFee + distanceKm*perKmRate) * multiplier

	order.SurgeMultiplier = multiplier
	order.DeliveryFee = math.Round(fee*100) / 100
}
