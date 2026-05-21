package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/FoPQer/go-shortener/internal/model"
	repo "github.com/FoPQer/go-shortener/internal/repository/user"
)

type MemoryUserRepository struct {
	users []*model.User
}

func NewRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make([]*model.User, 0),
	}
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	for _, u := range r.users {
		if u.GetID() == id {
			return u, nil
		}
	}
	return nil, repo.ErrUserNotFound
}

func (r *MemoryUserRepository) Save(ctx context.Context, user *model.User) (string, error) {
	if user.GetID() == "" {
		user.SetID(generateUserID())
	}
	for _, u := range r.users {
		if u.GetID() == user.GetID() {
			return "", fmt.Errorf("error while saving new user: %w", repo.ErrUserAlreadyExists)
		}
	}

	r.users = append(r.users, user)
	return user.GetID(), nil
}

func (r *MemoryUserRepository) GetUserURLs(ctx context.Context, userID string) ([]*model.Urls, error) {
	user, err := r.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.GetURLs(), nil
}

func generateUserID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}
