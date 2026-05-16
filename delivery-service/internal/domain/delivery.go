package domain

import "context"

type DeliveryStatus string

const (
	StatusPending   DeliveryStatus = "PENDING"
	StatusAssigned  DeliveryStatus = "ASSIGNED"
	StatusPickedUp  DeliveryStatus = "PICKED_UP"
	StatusDelivering DeliveryStatus = "DELIVERING"
	StatusCompleted  DeliveryStatus = "COMPLETED"
)

const (
	KafkaTopicPaymentProcessed = "payment-processed"
	KafkaGroupDelivery        = "delivery-group"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Courier struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Location      Location `json:"location"`
	IsAvailable   bool     `json:"is_available"`
	CurrentOrders []string `json:"current_orders"`
}

type Delivery struct {
	OrderID      string         `json:"order_id"`
	RestaurantID string         `json:"restaurant_id"`
	Status       DeliveryStatus `json:"status"`
	CourierID    string         `json:"courier_id,omitempty"`
	Destination  Location       `json:"destination"`
}

type CourierRepository interface {
	FindAll(ctx context.Context) ([]*Courier, error)
	Update(ctx context.Context, courier *Courier) error
	SaveCourier(ctx context.Context, courier *Courier) error // For seeding
}

type DeliveryRepository interface {
	Save(ctx context.Context, delivery *Delivery) error
	FindByID(ctx context.Context, orderID string) (*Delivery, error)
}
