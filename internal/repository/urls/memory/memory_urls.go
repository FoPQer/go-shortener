package memory

import (
	"context"
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

func (r *MemoryUrlsRepository) GetUrls(ctx context.Context) []*model.Urls {
	return r.urls
}

func (r *MemoryUrlsRepository) SetUrls(ctx context.Context, newUrls []*model.Urls) {
	r.urls = newUrls
}

func (r *MemoryUrlsRepository) GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error) {
	outUrls := make([]*model.Urls, 0)
	for _, u := range r.urls {
		if u.GetUserID() == userID && !u.IsDeleted() {
			outUrls = append(outUrls, u)
		}
	}
	return outUrls, fmt.Errorf("%w: %s", urls.ErrURLNotFound, userID)
}

func (r *MemoryUrlsRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (*model.Urls, error) {
	for _, u := range r.urls {
		if u.GetOriginal() == originalURL {
			if u.IsDeleted() {
				return nil, urls.ErrURLDeleted
			}
			return u, nil
		}
	}
	return nil, fmt.Errorf("error find by original URL %s: %w", originalURL, urls.ErrURLNotFound)
}

func (r *MemoryUrlsRepository) GetURLByShortURL(ctx context.Context, shortURL string) (string, error) {
	for _, u := range r.urls {
		if u.GetShortURL() == shortURL {
			if u.IsDeleted() {
				return "", urls.ErrURLDeleted
			}
			return u.GetOriginal(), nil
		}
	}
	return "", fmt.Errorf("error find by short URL %s: %w", shortURL, urls.ErrURLNotFound)
}

func (r *MemoryUrlsRepository) AddURL(ctx context.Context, original, shortURL string, userID string) (*model.Urls, error) {
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

func (r *MemoryUrlsRepository) AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error) {
	r.urls = append(r.urls, batchURLs...)
	return batchURLs, nil
}

func (r *MemoryUrlsRepository) DeleteUrls(ctx context.Context, shortUrls []string, userID string) error {
	for _, u := range r.urls {
		if slices.Contains(shortUrls, u.GetShortURL()) && u.GetUserID() == userID {
			u.SetDeleted(true)
		}
	}

	return nil
}
