package postgres

import (
	"context"
	"database/sql"
	"rinha_backend/internal/domain/models"
	"time"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (models.Client, error) {
	row := r.db.QueryRowContext(ctx, "SELECT account_limit, balance FROM clients WHERE id = $1", id)

	var accountLimit int64
	var balance int64
	err := row.Scan(&accountLimit, &balance)
	if err != nil {
		return models.Client{}, err
	}

	return models.Client{
		Date:         time.Now().Format(time.RFC3339Nano),
		AccountLimit: accountLimit,
		Balance:      balance,
	}, nil
}

func (r *ClientRepository) UpdateBalance(ctx context.Context, id string, newBalance int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE clients SET balance = $1 WHERE id = $2", newBalance, id)
	return err
}
