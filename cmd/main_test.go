package main

import (
	"audit-system/internal/database"
	"audit-system/internal/model"
	"audit-system/internal/repository"
	"audit-system/internal/service"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBankingSystem(t *testing.T) {
	database.Init()
	defer database.Close()

	userRepo := repository.NewUserRepository(database.Client)
	accountRepo := repository.NewAccountRepository(database.Client)
	transactionRepo := repository.NewTransactionRepository(database.Client)

	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	ctx := context.Background()

	user1 := model.User{
		Email: "user1@example.com",
		Name:  "User One",
		Age:   30,
	}
	err := userService.CreateUser(ctx, user1)
	assert.NoError(t, err, "failed to create user1")

	account1 := model.Account{
		Balance:          1000.0,
		LastTransferTime: time.Now(),
	}
	err = accountService.CreateAccount(ctx, account1, user1.Email)
	assert.NoError(t, err, "failed to create account for user1")

	user2 := model.User{
		Email: "user2@example.com",
		Name:  "User Two",
		Age:   25,
	}
	err = userService.CreateUser(ctx, user2)
	assert.NoError(t, err, "failed to create user2")

	account2 := model.Account{
		Balance:          500.0,
		LastTransferTime: time.Now(),
	}
	err = accountService.CreateAccount(ctx, account2, user2.Email)
	assert.NoError(t, err, "failed to create account for user2")

	accounts1, err := accountService.GetAccountsByEmail(ctx, user1.Email)
	assert.NoError(t, err, "failed to get accounts for user1")
	assert.NotEmpty(t, accounts1, "user1 should have one account")
	account1ID := accounts1[0].ID

	accounts2, err := accountService.GetAccountsByEmail(ctx, user2.Email)
	assert.NoError(t, err, "failed to get accounts for user2")
	assert.NotEmpty(t, accounts2, "user2 should have one account")
	account2ID := accounts2[0].ID

	transaction := model.Transaction{
		Amount:        200.0,
		Timestamp:     time.Now(),
		FromAccountID: account1ID,
		ToAccountID:   account2ID,
	}
	err = transactionService.CreateTransaction(ctx, user1.Email, account1ID, account2ID, transaction.Amount)
	assert.NoError(t, err, "failed to create transaction from user1 to user2")

	updatedAccount1, err := accountService.GetAccountByID(ctx, user1.Email, account1ID)
	assert.NoError(t, err, "failed to get updated account for user1")
	assert.Equal(t, 800.0, updatedAccount1.Balance, "balance of user1's account should be updated")

	updatedAccount2, err := accountService.GetAccountByID(ctx, user2.Email, account2ID)
	assert.NoError(t, err, "failed to get updated account for user2")
	assert.Equal(t, 700.0, updatedAccount2.Balance, "balance of user2's account should be updated")
}
