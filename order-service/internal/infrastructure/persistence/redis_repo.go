package persistence

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/user/highload-system-design/order-service/internal/domain"
)

type RedisCartRepository struct {
	client *redis.Client
}

func NewRedisCartRepository(client *redis.Client) *RedisCartRepository {
	return &RedisCartRepository{client: client}
}

func (r *RedisCartRepository) Save(ctx context.Context, c *domain.Cart) error {
	data, _ := json.Marshal(c)
	return r.client.Set(ctx, "cart:"+c.CustomerID, data, 7*24*time.Hour).Err()
}

func (r *RedisCartRepository) FindByCustomerID(ctx context.Context, customerID string) (*domain.Cart, error) {
	data, err := r.client.Get(ctx, "cart:"+customerID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var c domain.Cart
	json.Unmarshal(data, &c)
	return &c, nil
}

func (r *RedisCartRepository) DeleteByCustomerID(ctx context.Context, customerID string) error {
	return r.client.Del(ctx, "cart:"+customerID).Err()
}
