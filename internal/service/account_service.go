package service

import (
	"audit-system/internal/model"
	"audit-system/internal/repository"
	"context"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func newAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, account model.Account, userEmail string) (*model.Account, error) {
	return s.repo.CreateAccount(ctx, account, userEmail)
}

func (s *AccountService) GetAccountsByEmail(ctx context.Context, email string) ([]model.Account, error) {
	return s.repo.GetAccountsByEmail(ctx, email)
}
func (s *AccountService) GetAccountByID(ctx context.Context, email string, accountID int) (*model.Account, error) {
	return s.repo.GetAccountByID(ctx, email, accountID)
}
func (s *AccountService) GetAccountWithTransactions(ctx context.Context, email string, accountID int) (*model.Account, error) {
	return s.repo.GetAccountWithTransactions(ctx, email, accountID)
}

func (s *AccountService) UpdateAccount(ctx context.Context, email string, accountID int, updatedAccount model.Account) error {
	return s.repo.UpdateAccount(ctx, email, accountID, updatedAccount)
}
func (s *AccountService) DeleteAccount(ctx context.Context, accountID int) error {
	return s.repo.DeleteAccount(ctx, accountID)
}
