package service

import (
	"context"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/user"
)

type UserService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Get(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *model.User) (string, error) {
	return s.userRepo.Save(ctx, user)
}
