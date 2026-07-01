package memory

import (
	"context"
	"fmt"
	"slices"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
)

// MemoryUrlsRepository stores URL data in memory.
type MemoryUrlsRepository struct {
	urls []*model.Urls
}

// NewRepository creates an in-memory URL repository.
func NewRepository() *MemoryUrlsRepository {
	return &MemoryUrlsRepository{
		urls: make([]*model.Urls, 0),
	}
}

// GetUrls returns all URLs stored in memory.
func (r *MemoryUrlsRepository) GetUrls(ctx context.Context) []*model.Urls {
	return r.urls
}

// Count returns total amount of shortened URLs in memory.
func (r *MemoryUrlsRepository) Count(ctx context.Context) (int, error) {
	return len(r.urls), nil
}

// SetUrls replaces the in-memory URL collection.
func (r *MemoryUrlsRepository) SetUrls(ctx context.Context, newUrls []*model.Urls) {
	r.urls = newUrls
}

// GetUrlsByUserID returns non-deleted URLs that belong to the specified user.
func (r *MemoryUrlsRepository) GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error) {
	outUrls := make([]*model.Urls, 0)
	for _, u := range r.urls {
		if u.GetUserID() == userID && !u.IsDeleted() {
			outUrls = append(outUrls, u)
		}
	}
	return outUrls, nil
}

// GetURLByOriginalURL finds a URL entity by its original URL.
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

// GetURLByShortURL resolves a short URL token to its original URL.
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

// AddURL inserts a new URL into memory and returns the created entity.
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

// AddBatchURL appends a batch of URLs to memory and returns the stored entities.
func (r *MemoryUrlsRepository) AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error) {
	r.urls = append(r.urls, batchURLs...)
	return batchURLs, nil
}

// DeleteUrls marks matching URLs as deleted for the given user.
func (r *MemoryUrlsRepository) DeleteUrls(ctx context.Context, shortUrls []string, userID string) error {
	for _, u := range r.urls {
		if slices.Contains(shortUrls, u.GetShortURL()) && u.GetUserID() == userID {
			u.SetDeleted(true)
		}
	}

	return nil
}
