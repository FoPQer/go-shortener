package user

import (
	"context"
	"errors"

	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	Save(ctx context.Context, user *model.User) (string, error)
}