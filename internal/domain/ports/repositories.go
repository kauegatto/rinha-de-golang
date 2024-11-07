package ports

import (
	"context"
	"rinha_backend/internal/domain/models"
)

type ClientRepository interface {
	GetByID(ctx context.Context, id string) (models.Client, error)
	UpdateBalance(ctx context.Context, id string, newBalance int64) error
}

type TransactionRepository interface {
	Create(ctx context.Context, clientID string, transaction models.Transaction) error
	GetLastTenByClientID(ctx context.Context, clientID string) ([]models.Transaction, error)
}
