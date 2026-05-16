package persistence

import (
	"context"
	"sync"

	"github.com/user/highload-system-design/delivery-service/internal/domain"
)

type InMemoryRepository struct {
	mu         sync.RWMutex
	couriers   map[string]*domain.Courier
	deliveries map[string]*domain.Delivery
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		couriers:   make(map[string]*domain.Courier),
		deliveries: make(map[string]*domain.Delivery),
	}
}

// CourierRepository
func (r *InMemoryRepository) FindAll(ctx context.Context) ([]*domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Courier, 0, len(r.couriers))
	for _, c := range r.couriers {
		res = append(res, c)
	}
	return res, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, c *domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.couriers[c.ID] = c
	return nil
}

func (r *InMemoryRepository) SaveCourier(ctx context.Context, c *domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.couriers[c.ID] = c
	return nil
}

// DeliveryRepository
func (r *InMemoryRepository) Save(ctx context.Context, d *domain.Delivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries[d.OrderID] = d
	return nil
}

func (r *InMemoryRepository) FindByID(ctx context.Context, orderID string) (*domain.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.deliveries[orderID], nil
}
