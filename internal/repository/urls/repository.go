package urls

import "github.com/FoPQer/go-shortener/internal/model"

type Repository interface {
	GetUrls() []*model.Urls
	SetUrls(newUrls []*model.Urls)
	GetURLByShortURL(shortURL string) (string, error)
	AddURL(original, shortURL string) *model.Urls
}