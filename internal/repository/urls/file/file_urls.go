package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/FoPQer/go-shortener/internal/model"
	repository "github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/FoPQer/go-shortener/internal/utils"
)

// FileUrlsRepository stores URL data in a JSON file.
type FileUrlsRepository struct {
	filePath string
}

// NewRepository creates a file-backed URL repository.
func NewRepository(filePath string) *FileUrlsRepository {
	return &FileUrlsRepository{
		filePath: filePath,
	}
}

// GetUrls returns all URLs stored in the file.
func (r *FileUrlsRepository) GetUrls(ctx context.Context) []*model.Urls {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return make([]*model.Urls, 0)
	}

	var urls []*model.Urls
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return make([]*model.Urls, 0)
	}

	return urls
}

// SetUrls replaces the file contents with the provided URL collection.
func (r *FileUrlsRepository) SetUrls(ctx context.Context, newUrls []*model.Urls) {
	utils.WriteToFile(r.filePath, newUrls)
}

// GetURLByOriginalURL finds a URL entity by its original URL.
func (r *FileUrlsRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (*model.Urls, error) {
	urls := r.GetUrls(ctx)
	for _, u := range urls {
		if u.GetOriginal() == originalURL && !u.IsDeleted() {
			return u, nil
		}
	}
	return nil, fmt.Errorf("error find by original URL %s: %w", originalURL, repository.ErrURLNotFound)
}

// GetUrlsByUserID returns non-deleted URLs that belong to the specified user.
func (r *FileUrlsRepository) GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return make([]*model.Urls, 0), err
	}

	var urls []*model.Urls
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return make([]*model.Urls, 0), err
	}

	outUrls := make([]*model.Urls, 0)
	for _, u := range urls {
		if u.GetUserID() == userID && !u.IsDeleted() {
			outUrls = append(outUrls, u)
		}
	}

	return outUrls, nil
}

// GetURLByShortURL resolves a short URL token to its original URL.
func (r *FileUrlsRepository) GetURLByShortURL(ctx context.Context, shortURL string) (string, error) {
	urls := r.GetUrls(ctx)
	for _, u := range urls {
		if u.GetShortURL() == shortURL && !u.IsDeleted() {
			if u.IsDeleted() {
				return "", repository.ErrURLDeleted
			}
			return u.GetOriginal(), nil
		}
	}
	return "", fmt.Errorf("error find by short URL %s: %w", shortURL, repository.ErrURLNotFound)
}

// AddURL inserts a new URL into the file and returns the created entity.
func (r *FileUrlsRepository) AddURL(ctx context.Context, original, shortURL string, userID string) (*model.Urls, error) {
	urls := r.GetUrls(ctx)
	for _, u := range urls {
		if u.GetOriginal() == original && !u.IsDeleted() && u.GetUserID() == userID {
			return u, repository.ErrURLAlreadyExists
		}
	}

	u := model.NewUrls(original, shortURL)
	u.SetUserID(userID)

	urls = append(urls, u)

	utils.WriteToFile(r.filePath, urls)

	return u, nil
}

// AddBatchURL appends a batch of URLs to the file and returns the stored entities.
func (r *FileUrlsRepository) AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error) {
	urls := r.GetUrls(ctx)
	urls = append(urls, batchURLs...)
	utils.WriteToFile(r.filePath, urls)

	return batchURLs, nil
}

// DeleteUrls marks matching URLs as deleted and persists the updated file.
func (r *FileUrlsRepository) DeleteUrls(ctx context.Context, shortUrls []string, userID string) error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var urls []*model.Urls
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return err
	}

	for _, u := range urls {
		if slices.Contains(shortUrls, u.GetShortURL()) && u.GetUserID() == userID {
			u.SetDeleted(true)
		}
	}

	utils.WriteToFile(r.filePath, urls)

	return nil
}
