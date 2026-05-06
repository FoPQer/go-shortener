package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
)

type URLService struct {
	repo urls.Repository
}

func NewURLService(repo urls.Repository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) SetUrls(ctx context.Context, urls []*model.Urls) {
	s.repo.SetUrls(ctx, urls)
}

func (s *URLService) GetUrls(ctx context.Context) []*model.Urls {
	return s.repo.GetUrls(ctx)
}

func (s *URLService) GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error) {
	return s.repo.GetUrlsByUserID(ctx, userID)
}

func (s *URLService) DeleteUrls(ctx context.Context, shortUrls []string, userID string) error {
	if len(shortUrls) == 0 {
		return nil
	}

	numWorkers := min(len(shortUrls), 4)

	urlChan := make(chan string, len(shortUrls))
	errChan := make(chan error, numWorkers)

	go func() {
		for _, url := range shortUrls {
			urlChan <- url
		}
		close(urlChan)
	}()

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Go(func() {
			for shortURL := range urlChan {
				// Удаляем по одному URL'у
				err := s.repo.DeleteUrls(ctx, []string{shortURL}, userID)
				if err != nil {
					errChan <- err
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors while deleting URLs: %v", errs)
	}

	return nil
}

func (s *URLService) GetURL(ctx context.Context, shortURL string) (string, error) {
	url, err := s.repo.GetURLByShortURL(ctx, shortURL)
	if errors.Is(err, urls.ErrURLNotFound) {
		return "", fmt.Errorf("URL not found: %w", err)
	} else if errors.Is(err, urls.ErrURLDeleted) {
		return "", fmt.Errorf("URL is deleted: %w", err)
	} else if err != nil {
		return "", fmt.Errorf("unable to get URL: %w", err)
	}

	return url, nil
}

func (s *URLService) SetURL(ctx context.Context, fullURL string, userID string) (string, error) {
	id := newID()
	url, err := s.repo.AddURL(ctx, fullURL, id, userID)
	if errors.Is(err, urls.ErrURLAlreadyExists) {
		short, makeErr := MakeShortURL(url.GetShortURL())
		if makeErr != nil {
			return "", errors.Join(fmt.Errorf("unsuccessful URL creation: %w", err), makeErr)
		}

		return short, urls.ErrURLAlreadyExists
	} else if err != nil {
		return "", fmt.Errorf("unable to add URL: %w", err)
	}
	logger.GetSugar().Infof("Added URL: %s -> %s", url.GetOriginal(), url.GetShortURL()) 

	return MakeShortURL(url.GetShortURL())
}

func (s *URLService) SetBatchURL(ctx context.Context, batchURLs []*model.Urls, userID string) ([]*model.Urls, error) {	
	for _, u := range batchURLs {
		u.SetUserID(string(userID))
	}
	result, err := s.repo.AddBatchURL(ctx, batchURLs)
	if err != nil {
		return nil, fmt.Errorf("unable to add batch URLs: %w", err)
	}

	return result, nil

}

func newID() string {
	return rand.Text()[0:8]
}

func MakeShortURL(id string) (string, error) {
	short, err := url.JoinPath("http://"+GetRunAddr(), GetBasePrefix(), id)
	if err != nil {
		return "", err
	}
	return short, nil
}
