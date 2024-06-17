package repository

import (
	"audit-system/ent"
	"audit-system/ent/account"
	"audit-system/ent/transaction"
	"audit-system/ent/user"
	"audit-system/internal/model"
	"context"
	"fmt"
	"time"
)

type TransactionRepository struct {
	client *ent.Client
}

func NewTransactionRepository(client *ent.Client) *TransactionRepository {
	return &TransactionRepository{client: client}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, email string, fromAccountID int, toAccountID int, amount float64) error {
	fromAccount, err := r.client.Account.Query().
		Where(account.IDEQ(fromAccountID), account.HasUserWith(user.EmailEQ(email))).
		Only(ctx)
	if err != nil {
		return err
	}
	toAccount, err := r.client.Account.Get(ctx, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to get toAccount: %w", err)
	}

	// Update balances and last transfer times
	fromAccountUpdate := r.client.Account.UpdateOneID(fromAccountID).
		AddBalance(-amount).
		SetLastTransferTime(time.Now())

	toAccountUpdate := r.client.Account.UpdateOneID(toAccountID).
		AddBalance(amount).
		SetLastTransferTime(time.Now())

	_, err = r.client.Transaction.
		Create().
		SetAmount(amount).
		SetTimestamp(time.Now()).
		SetFromAccount(fromAccount).
		SetToAccount(toAccount).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Save the account updates
	if _, err = fromAccountUpdate.Save(ctx); err != nil {
		return fmt.Errorf("failed to update fromAccount: %w", err)
	}
	if _, err = toAccountUpdate.Save(ctx); err != nil {
		return fmt.Errorf("failed to update toAccount: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetTransactions(ctx context.Context, email string, accountID int, transactionType string) ([]model.Transaction, error) {
	var transactions []*ent.Transaction
	var err error

	switch transactionType {
	case "inbound":
		transactions, err = r.client.Transaction.Query().
			Where(transaction.HasToAccountWith(account.IDEQ(accountID), account.HasUserWith(user.EmailEQ(email)))).
			WithFromAccount().
			WithToAccount().
			All(ctx)
	case "outbound":
		transactions, err = r.client.Transaction.Query().
			Where(transaction.HasFromAccountWith(account.IDEQ(accountID), account.HasUserWith(user.EmailEQ(email)))).
			WithFromAccount().
			WithToAccount().
			All(ctx)
	default:
		transactions, err = r.client.Transaction.Query().
			Where(transaction.Or(
				transaction.HasFromAccountWith(account.IDEQ(accountID), account.HasUserWith(user.EmailEQ(email))),
				transaction.HasToAccountWith(account.IDEQ(accountID), account.HasUserWith(user.EmailEQ(email))),
			)).
			WithFromAccount().
			WithToAccount().
			All(ctx)
	}

	if err != nil {
		return nil, err
	}

	result := make([]model.Transaction, len(transactions))
	for i, t := range transactions {
		result[i] = model.Transaction{
			ID:            t.ID,
			Amount:        t.Amount,
			Timestamp:     t.Timestamp,
			FromAccountID: t.Edges.FromAccount.ID,
			ToAccountID:   t.Edges.ToAccount.ID,
		}
	}
	return result, nil
}
