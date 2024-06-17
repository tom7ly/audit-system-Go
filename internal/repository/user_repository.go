package repository

import (
	"audit-system/ent"
	"audit-system/ent/user"
	"audit-system/internal/model"
	"context"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.client.User.
		Create().
		SetEmail(user.Email).
		SetName(user.Name).
		SetAge(user.Age).
		Save(ctx)
	return err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := r.client.User.Query().
		WithAccounts().
		All(ctx)
	if err != nil {
		return nil, err
	}

	return mapUsers(users), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	u, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		WithAccounts().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	mappedUser := mapUser(u)
	return &mappedUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, email string, updatedUser model.User) error {
	_, err := r.client.User.
		Update().
		Where(user.EmailEQ(email)).
		SetName(updatedUser.Name).
		SetAge(updatedUser.Age).
		Save(ctx)
	return err
}

func mapUsers(users []*ent.User) []model.User {
	result := make([]model.User, len(users))
	for i, u := range users {
		result[i] = mapUser(u)
	}
	return result
}

func mapUser(u *ent.User) model.User {
	accounts := make([]model.Account, len(u.Edges.Accounts))
	for i, acc := range u.Edges.Accounts {
		accounts[i] = model.Account{
			ID:               acc.ID,
			Balance:          acc.Balance,
			LastTransferTime: acc.LastTransferTime,
		}
	}

	return model.User{
		Email:    u.Email,
		Name:     u.Name,
		Age:      u.Age,
		Accounts: accounts,
	}
}
