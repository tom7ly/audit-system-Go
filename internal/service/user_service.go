package service

import (
	"audit-system/internal/model"
	"audit-system/internal/repository"
	"context"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
func (s *UserService) UpdateUser(ctx context.Context, email string, updatedUser model.User) error {
	return s.repo.UpdateUser(ctx, email, updatedUser)
}
