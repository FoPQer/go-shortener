package user

import (
	"context"
	"errors"

	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	// ErrUserNotFound indicates that a user record does not exist.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists indicates that a user with the same identity already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository defines persistence operations for user entities.
type UserRepository interface {
	// FindByID returns a user by its identifier.
	FindByID(ctx context.Context, id string) (*model.User, error)
	// Count returns total amount of users.
	Count(ctx context.Context) (int, error)
	// Save persists a user and returns its identifier.
	Save(ctx context.Context, user *model.User) (string, error)
}
