package service

import (
	"audit-system/internal/model"
	"audit-system/internal/repository"
	"context"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, email string, fromAccountID int, toAccountID int, amount float64) error {
	return s.repo.CreateTransaction(ctx, email, fromAccountID, toAccountID, amount)
}

func (s *TransactionService) GetTransactions(ctx context.Context, email string, accountID int, transactionType string) ([]model.Transaction, error) {
	return s.repo.GetTransactions(ctx, email, accountID, transactionType)
}
