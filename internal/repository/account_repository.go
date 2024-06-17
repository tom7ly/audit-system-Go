package repository

import (
	"audit-system/ent"
	"audit-system/ent/account"
	"audit-system/ent/user"
	"audit-system/internal/model"
	"context"
)

type AccountRepository struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) *AccountRepository {
	return &AccountRepository{client: client}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account model.Account, userEmail string) error {
	userEntity, err := r.client.User.Query().Where(user.EmailEQ(userEmail)).Only(ctx)
	if err != nil {
		return err
	}

	_, err = r.client.Account.
		Create().
		SetBalance(account.Balance).
		SetLastTransferTime(account.LastTransferTime).
		SetUser(userEntity).
		Save(ctx)
	return err
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
	_, err := r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return err
	}
	_, err = r.client.Account.UpdateOneID(accountID).
		SetBalance(updatedAccount.Balance).
		Save(ctx)
	return err
}

func mapAccounts(accounts []*ent.Account) []model.Account {
	result := make([]model.Account, len(accounts))
	for i, a := range accounts {
		result[i] = *mapAccount(a)
	}
	return result
}

func mapAccount(a *ent.Account) *model.Account {
	// Map transactions
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
		UserEmail:            a.Edges.User.Email,
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
