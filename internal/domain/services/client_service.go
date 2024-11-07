package services

import (
	"context"
	"rinha_backend/internal/domain/models"
	"rinha_backend/internal/domain/ports"
)

type ClientService struct {
	clientRepo      ports.ClientRepository
	transactionRepo ports.TransactionRepository
}

func NewClientService(cr ports.ClientRepository, tr ports.TransactionRepository) *ClientService {
	return &ClientService{
		clientRepo:      cr,
		transactionRepo: tr,
	}
}

func (s *ClientService) ProcessTransaction(ctx context.Context, clientID string, payment models.Payment) (models.Client, error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return models.Client{}, err
	}

	newBalance, err := payment.ValidateAndReturnNewBalance(client.Balance)
	if err != nil {
		return models.Client{}, err
	}

	transaction := models.Transaction{
		Amount:      payment.Value,
		Operation:   payment.Type,
		Description: payment.Description,
	}

	if err := s.transactionRepo.Create(ctx, clientID, transaction); err != nil {
		return models.Client{}, err
	}

	if err := s.clientRepo.UpdateBalance(ctx, clientID, newBalance); err != nil {
		return models.Client{}, err
	}

	client.Balance = newBalance
	return client, nil
}

func (s *ClientService) GetExtract(ctx context.Context, clientID string) (models.Client, []models.Transaction, error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return models.Client{}, nil, err
	}

	transactions, err := s.transactionRepo.GetLastTenByClientID(ctx, clientID)
	if err != nil {
		return models.Client{}, nil, err
	}

	return client, transactions, nil
}
