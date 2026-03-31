package urls

import (
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
	GetUrls() []*model.Urls
	SetUrls(newUrls []*model.Urls)
	GetUrlsByUserID(userID string) ([]*model.Urls, error)
	GetURLByOriginalURL(originalURL string) (*model.Urls, error)
	GetURLByShortURL(shortURL string) (string, error)
	AddURL(original, shortURL string, userID string) (*model.Urls, error)
	AddBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error)
	DeleteUrls(shortUrls []string, userID string) error
}