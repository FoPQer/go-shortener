package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/model"
	repo "github.com/FoPQer/go-shortener/internal/repository/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBUserRepository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *DBUserRepository {
	return &DBUserRepository{
		conn: conn,
	}
}

func (r *DBUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	row := r.conn.QueryRow(
		ctx, 
		"SELECT id FROM users WHERE id = $1",
		id,
	)

	err := row.Scan(&user.ID)
	logger.GetSugar().Infof("Finded user: %v", user)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repo.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}

func (r *DBUserRepository) Save(ctx context.Context, user *model.User) (string, error) {
	row := r.conn.QueryRow(
		ctx, 
		"INSERT INTO users DEFAULT VALUES RETURNING id",
	)
	err := row.Scan(&user.ID)
	
	if err != nil {
		return "", fmt.Errorf("failed to save user: %w", err)
	}
	return user.ID, nil
}

