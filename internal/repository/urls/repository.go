package urls

import (
	"errors"

	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	ErrBadValueReceive  = errors.New("value not received")
	ErrURLAlreadyExists = errors.New("url is already exists")
)

type Repository interface {
	GetUrls() []*model.Urls
	SetUrls(newUrls []*model.Urls)
	GetURLByOriginalURL(originalURL string) (*model.Urls, error)
	GetURLByShortURL(shortURL string) (string, error)
	AddURL(original, shortURL string) (*model.Urls, error)
}