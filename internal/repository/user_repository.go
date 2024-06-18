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

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client, q *utils.Queue) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {

	u, err := r.client.User.Create().
		SetEmail(user.Email).
		SetName(user.Name).
		SetAge(user.Age).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	userModel := mapUser(u)
	return &userModel, nil
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

func (r *UserRepository) UpdateUser(ctx context.Context, email string, updatedUser model.User) (*model.User, error) {
	var updatedEntUser *ent.User
	err := utils.WithTx(r.client, ctx, func(tx *ent.Tx) error {
		userEntity, err := tx.User.Query().Where(user.EmailEQ(email)).Only(ctx)
		if err != nil {
			return err
		}

		updatedEntUser, err = tx.User.UpdateOne(userEntity).
			SetName(updatedUser.Name).
			SetAge(updatedUser.Age).
			Save(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	updatedModelUser := mapUser(updatedEntUser)
	return &updatedModelUser, nil

}
func (r *UserRepository) DeleteUser(ctx context.Context, email string) error {
	return utils.WithTx(r.client, ctx, func(tx *ent.Tx) error {
		userEntity, err := tx.User.Query().Where(user.EmailEQ(email)).Only(ctx)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}

		_, err = tx.Account.Delete().Where(account.HasUserWith(user.EmailEQ(email))).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete accounts for user: %w", err)
		}

		err = tx.User.DeleteOne(userEntity).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		return nil
	})
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
