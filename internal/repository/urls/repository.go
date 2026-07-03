package urls

import (
	"context"
	"errors"

	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	// ErrBadValueReceive indicates that required input data was not provided.
	ErrBadValueReceive = errors.New("value not received")
	// ErrURLAlreadyExists indicates that a URL with the same original value already exists.
	ErrURLAlreadyExists = errors.New("url is already exists")
	// ErrURLNotFound indicates that a URL record was not found.
	ErrURLNotFound = errors.New("url not found")
	// ErrURLDeleted indicates that the requested URL exists but was marked as deleted.
	ErrURLDeleted = errors.New("url is deleted")
)

// Repository defines storage operations for shortened URLs.
//
// Implementations may use in-memory, file-based, or database-backed storage.
type Repository interface {
	// GetUrls returns all stored URLs.
	GetUrls(ctx context.Context) []*model.Urls
	// Count returns total amount of shortened URLs.
	Count(ctx context.Context) (int, error)
	// SetUrls replaces or initializes the storage with the provided URLs.
	SetUrls(ctx context.Context, newUrls []*model.Urls)
	// GetUrlsByUserID returns all non-deleted URLs that belong to a user.
	GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error)
	// GetURLByOriginalURL finds a URL by its original value.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (*model.Urls, error)
	// GetURLByShortURL resolves a short URL to its original URL string.
	GetURLByShortURL(ctx context.Context, shortURL string) (string, error)
	// AddURL stores a new URL for a user and returns the created entity.
	AddURL(ctx context.Context, original, shortURL string, userID string) (*model.Urls, error)
	// AddBatchURL stores multiple URLs in a single operation.
	AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error)
	// DeleteUrls marks the provided short URLs as deleted for the given user.
	DeleteUrls(ctx context.Context, shortUrls []string, userID string) error
}
