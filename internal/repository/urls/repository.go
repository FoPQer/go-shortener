package urls

import (
	"context"
	"errors"

	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	ErrBadValueReceive  = errors.New("value not received")
	ErrURLAlreadyExists = errors.New("url is already exists")
	ErrURLNotFound      = errors.New("url not found")
	ErrURLDeleted       = errors.New("url is deleted")
)

type Repository interface {
	GetUrls(ctx context.Context) []*model.Urls
	SetUrls(ctx context.Context, newUrls []*model.Urls)
	GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (*model.Urls, error)
	GetURLByShortURL(ctx context.Context, shortURL string) (string, error)
	AddURL(ctx context.Context, original, shortURL string, userID string) (*model.Urls, error)
	AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error)
	DeleteUrls(ctx context.Context, shortUrls []string, userID string) error
}