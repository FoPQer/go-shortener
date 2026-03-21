package service

import (
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

func (s *UserService) Get(id string) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Create(user *model.User) (string, error) {
	return s.userRepo.Save(user)
}