package persistence

import (
	"context"
	"database/sql"
	"github.com/user/highload-system-design/payment-service/internal/domain"
)

type PostgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Save(ctx context.Context, p *domain.Payment) error {
	_, err := r.db.ExecContext(ctx, 
		"INSERT INTO payments (id, order_id, amount, status, idempotency_key, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		p.ID, p.OrderID, p.Amount, p.Status, p.IdempotencyKey, p.CreatedAt)
	return err
}

func (r *PostgresPaymentRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, order_id, amount, status, idempotency_key, created_at FROM payments WHERE idempotency_key = $1", key)
	p := &domain.Payment{}
	err := row.Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.IdempotencyKey, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}
