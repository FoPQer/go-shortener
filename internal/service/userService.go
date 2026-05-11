package service

import (
	"context"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/user"
)

// UserService provides business operations for user retrieval and creation.
type UserService struct {
	userRepo user.UserRepository
}

// NewUserService constructs a UserService with the given user repository.
func NewUserService(userRepo user.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Get returns a user by its identifier.
func (s *UserService) Get(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// Create persists a new user and returns its identifier.
func (s *UserService) Create(ctx context.Context, user *model.User) (string, error) {
	return s.userRepo.Save(ctx, user)
}
