package domain

import (
	"context"
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusDelivering OrderStatus = "DELIVERING"
	StatusCompleted  OrderStatus = "COMPLETED"
)

type OrderItem struct {
	MenuItemID string  `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type Order struct {
	OrderID           string      `json:"order_id"`
	CustomerID        string      `json:"customer_id"`
	RestaurantID      string      `json:"restaurant_id"`
	Status            OrderStatus `json:"status"`
	TotalAmount       float64     `json:"total_amount"`
	DeliveryFee       float64     `json:"delivery_fee"`
	SurgeMultiplier   float64     `json:"surge_multiplier"`
	Items             []OrderItem `json:"items"`
	DeliveryAddress   string      `json:"delivery_address"`
	EstimatedDelivery time.Time   `json:"estimated_delivery"`
	CreatedAt         time.Time   `json:"created_at"`
}

type CartItem struct {
	MenuItemID string `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
}

type Cart struct {
	CustomerID   string     `json:"customer_id"`
	RestaurantID string     `json:"restaurant_id"`
	Items        []CartItem `json:"items"`
}

type OutboxEvent struct {
	ID            string    `json:"id"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	EventType     string    `json:"event_type"`
	Payload       string    `json:"payload"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
}

type CartRepository interface {
	Save(ctx context.Context, cart *Cart) error
	FindByCustomerID(ctx context.Context, customerID string) (*Cart, error)
	DeleteByCustomerID(ctx context.Context, customerID string) error
}

type OutboxRepository interface {
	SaveEvent(ctx context.Context, event *OutboxEvent) error
	FindPending(ctx context.Context) ([]OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, id string) error
}
