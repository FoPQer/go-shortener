package memory

import (
	"fmt"
	"slices"

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

func (r *MemoryUrlsRepository) GetUrlsByUserID(userID string) ([]*model.Urls, error) {
	out_urls := make([]*model.Urls, 0)
	for _, u := range r.urls {
		if u.GetUserID() == userID {
			out_urls = append(out_urls, u)
		}
	}
	return out_urls, fmt.Errorf("%w: %s", urls.ErrUrlNotFound, userID)
}

func (r *MemoryUrlsRepository) GetURLByOriginalURL(originalURL string) (*model.Urls, error) {
	for _, u := range r.urls {
		if u.GetOriginal() == originalURL {
			return u, nil
		}
	}
	return nil, fmt.Errorf("error find by original URL %s: %w", originalURL, urls.ErrUrlNotFound)
}


func (r *MemoryUrlsRepository) GetURLByShortURL(shortURL string) (string, error) {
	for _, u := range r.urls {
		if u.GetShortURL() == shortURL {
			return u.GetOriginal(), nil
		}
	}
	return "", fmt.Errorf("error find by short URL %s: %w", shortURL, urls.ErrUrlNotFound)
}

func (r *MemoryUrlsRepository) AddURL(original, shortURL, userID string) (*model.Urls, error) {
	for _, u := range r.urls {
		if u.GetOriginal() == original {
			return u, urls.ErrURLAlreadyExists
		}
	}

	u := model.NewUrls(original, shortURL)
	u.SetUserID(userID)
	r.urls = append(r.urls, u)
	return u, nil
}

func (r *MemoryUrlsRepository) AddBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {
	r.urls = append(r.urls, batchURLs...)
	return batchURLs, nil
}

func (r *MemoryUrlsRepository) DeleteUrlsByUserID(userID string) error {
	for i, u := range r.urls {
		if u.GetUserID() == userID {
			r.urls = slices.Delete(r.urls, i, i+1)
		}
	}

	return nil
}
