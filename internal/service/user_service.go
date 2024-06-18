package service

import (
	"audit-system/internal/model"
	"audit-system/internal/repository"
	"context"
)

type UserService struct {
	userRepo    *repository.UserRepository
	accountRepo *repository.AccountRepository
}

func newUserService(userRepo *repository.UserRepository, accountRepo *repository.AccountRepository) *UserService {
	return &UserService{userRepo: userRepo, accountRepo: accountRepo}
}
func (s *UserService) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetUsers(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}
func (s *UserService) UpdateUser(ctx context.Context, email string, updatedUser model.User) (*model.User, error) {
	return s.userRepo.UpdateUser(ctx, email, updatedUser)
}
func (s *UserService) DeleteUser(ctx context.Context, email string) error {
	return s.userRepo.DeleteUser(ctx, email)
}
