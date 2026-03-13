package memory

import (
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
)

type MemoryUrlsRepository struct {
	urls []*model.Urls
}

func NewRepository() *MemoryUrlsRepository {
	return &MemoryUrlsRepository{
		urls: make([]*model.Urls, 0),
	}
}

func (r *MemoryUrlsRepository) GetUrls() []*model.Urls {
	return r.urls
}

func (r *MemoryUrlsRepository) SetUrls(newUrls []*model.Urls) {
	r.urls = newUrls
}

func (r *MemoryUrlsRepository) GetURLByOriginalURL(originalURL string) (*model.Urls, error) {
	for _, u := range r.urls {
		if u.GetOriginal() == originalURL {
			return u, nil
		}
	}
	return nil, urls.ErrBadValueReceive
}


func (r *MemoryUrlsRepository) GetURLByShortURL(shortURL string) (string, error) {
	for _, u := range r.urls {
		if u.GetShortURL() == shortURL {
			return u.GetOriginal(), nil
		}
	}
	return "", urls.ErrBadValueReceive
}

func (r *MemoryUrlsRepository) AddURL(original, shortURL string) (*model.Urls, error) {
	u := model.NewUrls(original, shortURL)
	r.urls = append(r.urls, u)
	return u, nil
}
