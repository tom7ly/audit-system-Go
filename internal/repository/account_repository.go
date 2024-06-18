package repository

import (
	"audit-system/ent"
	"audit-system/ent/account"
	"audit-system/ent/user"
	"audit-system/internal/model"
	"audit-system/internal/utils"
	"context"
	"fmt"
)

type AccountRepository struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client, q *utils.Queue) *AccountRepository {
	return &AccountRepository{client: client}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, newAccount model.Account, userEmail string) (*model.Account, error) {
	var createdEntAccount *ent.Account
	err := utils.WithTx(r.client, ctx, func(tx *ent.Tx) error {
		userEntity, err := tx.User.Query().Where(user.EmailEQ(userEmail)).Only(ctx)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}

		createdEntAccount, err = tx.Account.
			Create().
			SetBalance(newAccount.Balance).
			SetLastTransferTime(newAccount.LastTransferTime).
			SetUser(userEntity).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		// Reload the account with the edges to ensure the user edge is set
		createdEntAccount, err = tx.Account.Query().
			Where(account.IDEQ(createdEntAccount.ID)).
			WithUser().
			Only(ctx)
		if err != nil {
			return fmt.Errorf("failed to reload account: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapAccount(createdEntAccount), nil
}

func (r *AccountRepository) GetAccountsByEmail(ctx context.Context, email string) ([]model.Account, error) {
	accounts, err := r.client.Account.Query().
		Where(account.HasUserWith(user.EmailEQ(email))).WithUser().
		WithOutgoingTransactions(func(q *ent.TransactionQuery) {
			q.WithToAccount()
		}).
		WithIncomingTransactions(func(q *ent.TransactionQuery) {
			q.WithFromAccount()
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return mapAccounts(accounts), nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, email string, accountID int) (*model.Account, error) {
	accountEntity, err := r.client.Account.Query().
		Where(account.IDEQ(accountID), account.HasUserWith(user.EmailEQ(email))).
		WithIncomingTransactions(func(q *ent.TransactionQuery) {
			q.WithFromAccount()
		}).
		WithOutgoingTransactions(func(q *ent.TransactionQuery) {
			q.WithToAccount()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return mapAccount(accountEntity), nil
}

func (r *AccountRepository) GetAccountWithTransactions(ctx context.Context, email string, accountID int) (*model.Account, error) {
	accountEntity, err := r.client.Account.Query().
		Where(account.IDEQ(accountID), account.HasUserWith(user.EmailEQ(email))).WithUser().
		WithOutgoingTransactions(func(q *ent.TransactionQuery) {
			q.WithToAccount()
		}).
		WithIncomingTransactions(func(q *ent.TransactionQuery) {
			q.WithFromAccount()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return mapAccount(accountEntity), nil
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, email string, accountID int, updatedAccount model.Account) error {
	return utils.WithTx(r.client, ctx, func(tx *ent.Tx) error {
		_, err := tx.User.Query().Where(user.EmailEQ(email)).Only(ctx)
		if err != nil {
			return err
		}
		acc, err := tx.Account.Get(ctx, accountID)
		if err != nil {
			return err
		}

		_, err = tx.Account.UpdateOne(acc).
			SetBalance(updatedAccount.Balance).
			Save(ctx)
		return err
	})
}
func (r *AccountRepository) DeleteAccount(ctx context.Context, accountID int) error {
	return utils.WithTx(r.client, ctx, func(tx *ent.Tx) error {
		accountEntity, err := tx.Account.Get(ctx, accountID)
		if err != nil {
			return fmt.Errorf("failed to find account: %w", err)
		}

		err = tx.Account.DeleteOne(accountEntity).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete account: %w", err)
		}

		return nil
	})
}

func mapAccounts(accounts []*ent.Account) []model.Account {
	result := make([]model.Account, len(accounts))
	for i, a := range accounts {
		result[i] = *mapAccount(a)
	}
	return result
}

func mapAccount(a *ent.Account) *model.Account {
	outgoingTransactions := make([]model.Transaction, len(a.Edges.OutgoingTransactions))
	for i, t := range a.Edges.OutgoingTransactions {
		outgoingTransactions[i] = model.Transaction{
			ID:            t.ID,
			Amount:        t.Amount,
			Timestamp:     t.Timestamp,
			FromAccountID: safeAccountID(t.Edges.FromAccount),
			ToAccountID:   safeAccountID(t.Edges.ToAccount),
		}
	}

	incomingTransactions := make([]model.Transaction, len(a.Edges.IncomingTransactions))
	for i, t := range a.Edges.IncomingTransactions {
		incomingTransactions[i] = model.Transaction{
			ID:            t.ID,
			Amount:        t.Amount,
			Timestamp:     t.Timestamp,
			FromAccountID: safeAccountID(t.Edges.FromAccount),
			ToAccountID:   safeAccountID(t.Edges.ToAccount),
		}
	}

	return &model.Account{
		ID:                   a.ID,
		Balance:              a.Balance,
		LastTransferTime:     a.LastTransferTime,
		OutgoingTransactions: outgoingTransactions,
		IncomingTransactions: incomingTransactions,
	}
}

func safeAccountID(acc *ent.Account) int {
	if acc != nil {
		return acc.ID
	}
	return 0
}
