package services

import (
	"context"
	"fmt"
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

	errChan := make(chan error, 2)
	go func() {
		transaction := models.Transaction{
			Amount:      payment.Value,
			Operation:   payment.Type,
			Description: payment.Description,
		}

		errChan <- s.transactionRepo.Create(ctx, clientID, transaction)
	}()

	go func() {
		errChan <- s.clientRepo.UpdateBalance(ctx, clientID, newBalance)
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return models.Client{}, fmt.Errorf("parallel operation failed: %w", err)
		}
	}

	client.Balance = newBalance
	return client, nil
}

func (s *ClientService) GetExtract(ctx context.Context, clientID string) (models.Client, []models.Transaction, error) {
	type clientResult struct {
		client models.Client
		err    error
	}
	type transactionsResult struct {
		transactions []models.Transaction
		err          error
	}

	clientCh := make(chan clientResult, 1)
	transactionsCh := make(chan transactionsResult, 1)

	go func() {
		client, err := s.clientRepo.GetByID(ctx, clientID)
		clientCh <- clientResult{client, err}
	}()

	go func() {
		transactions, err := s.transactionRepo.GetLastTenByClientID(ctx, clientID)
		transactionsCh <- transactionsResult{transactions, err}
	}()

	clientRes := <-clientCh
	if clientRes.err != nil {
		return models.Client{}, nil, fmt.Errorf("getting client: %w", clientRes.err)
	}

	transactionsRes := <-transactionsCh
	if transactionsRes.err != nil {
		return models.Client{}, nil, fmt.Errorf("getting transactions: %w", transactionsRes.err)
	}

	return clientRes.client, transactionsRes.transactions, nil
}
