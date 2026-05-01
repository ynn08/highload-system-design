package persistence

import (
	"context"
	"database/sql"

	"github.com/user/highload-system-design/order-service/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// OrderRepository implementation
func (r *PostgresRepository) Save(ctx context.Context, o *domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO orders (order_id, customer_id, restaurant_id, status, total_amount, delivery_fee, surge_multiplier, delivery_address, estimated_delivery, created_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		o.OrderID, o.CustomerID, o.RestaurantID, o.Status, o.TotalAmount, o.DeliveryFee, o.SurgeMultiplier, o.DeliveryAddress, o.EstimatedDelivery, o.CreatedAt)
	return err
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	row := r.db.QueryRowContext(ctx, "SELECT order_id, customer_id, restaurant_id, status, total_amount, delivery_fee, surge_multiplier, delivery_address, estimated_delivery, created_at FROM orders WHERE order_id = $1", id)
	o := &domain.Order{}
	err := row.Scan(&o.OrderID, &o.CustomerID, &o.RestaurantID, &o.Status, &o.TotalAmount, &o.DeliveryFee, &o.SurgeMultiplier, &o.DeliveryAddress, &o.EstimatedDelivery, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// OutboxRepository implementation
// Note: Save for Outbox and Order are different in domain. 
// This struct will implement both. To avoid conflict if domain had both as "Save", 
// we ensure the method signature matches the interface.

func (r *PostgresRepository) SaveEvent(ctx context.Context, e *domain.OutboxEvent) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO outbox_events (id, aggregate_id, aggregate_type, event_type, payload, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		e.ID, e.AggregateID, e.AggregateType, e.EventType, e.Payload, e.Status, e.CreatedAt)
	return err
}

func (r *PostgresRepository) FindPending(ctx context.Context) ([]domain.OutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, aggregate_id, aggregate_type, event_type, payload, status, created_at FROM outbox_events WHERE status = 'PENDING'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var e domain.OutboxEvent
		if err := rows.Scan(&e.ID, &e.AggregateID, &e.AggregateType, &e.EventType, &e.Payload, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *PostgresRepository) MarkAsProcessed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE outbox_events SET status = 'PROCESSED' WHERE id = $1", id)
	return err
}
