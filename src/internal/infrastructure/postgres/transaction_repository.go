package postgres

import (
	"context"
	"rinha_backend/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, clientID string, transaction models.Transaction) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO transactions (client_id, amount, operation, description) VALUES ($1, $2, $3, $4)",
		clientID, transaction.Amount, transaction.Operation, transaction.Description)
	return err
}

func (r *TransactionRepository) GetLastTenByClientID(ctx context.Context, clientID string) ([]models.Transaction, error) {
	rows, err := r.db.Query(ctx,
		"SELECT amount, operation, description, created_at FROM transactions WHERE client_id = $1 ORDER BY created_at DESC LIMIT 10",
		clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.Amount, &t.Operation, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
