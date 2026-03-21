package user

import "github.com/FoPQer/go-shortener/internal/model"

type UserRepository interface {
	FindByID(id string) (*model.User, error)
	Save(user *model.User) (string, error)
}