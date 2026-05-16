package usecase

import (
	"context"
	"github.com/user/highload-system-design/payment-service/internal/domain"
	"testing"
)

type mockRepo struct {
	payments map[string]*domain.Payment
}

func (m *mockRepo) Save(ctx context.Context, p *domain.Payment) error {
	m.payments[p.IdempotencyKey] = p
	return nil
}

func (m *mockRepo) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error) {
	return m.payments[key], nil
}

type mockPublisher struct{}

func (m *mockPublisher) PublishPaymentProcessed(ctx context.Context, orderID string, restaurantID string, status string) error {
	return nil
}

type mockGateway struct {
	success bool
}

func (m *mockGateway) Charge(ctx context.Context, amount float64) (bool, error) {
	return m.success, nil
}

func TestProcessPaymentUseCase_Idempotency(t *testing.T) {
	repo := &mockRepo{payments: make(map[string]*domain.Payment)}
	pub := &mockPublisher{}
	gw := &mockGateway{success: true}
	uc := NewProcessPaymentUseCase(repo, pub, gw)

	orderID := "order-1"
	restaurantID := "rest-1"
	
	// First call
	err := uc.Execute(context.Background(), orderID, restaurantID, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	if len(repo.payments) != 1 {
		t.Errorf("expected 1 payment, got %d", len(repo.payments))
	}

	// Second call with same orderID (idempotency key)
	err = uc.Execute(context.Background(), orderID, restaurantID, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	if len(repo.payments) != 1 {
		t.Errorf("expected still 1 payment (idempotency), got %d", len(repo.payments))
	}
}
