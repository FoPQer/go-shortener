package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/url"

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

func (s *URLService) SetUrls(urls []*model.Urls) {
	s.repo.SetUrls(urls)
}

func (s *URLService) GetUrls() []*model.Urls {
	return s.repo.GetUrls()
}

func (s *URLService) GetURL(shortURL string) (string, error) {
	url, err := s.repo.GetURLByShortURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("Unable to get URL: %w", err)
	}

	return url, nil
}

func (s *URLService) SetURL(fullURL string) (string, error) {
	id := newID()
	url, err := s.repo.AddURL(fullURL, id)
	if errors.Is(err, urls.ErrURLAlreadyExists) {
		short, makeErr := makeShortURL(url.GetShortURL())
		if makeErr != nil {
			return "", errors.Join(fmt.Errorf("Unsuccessful URL creation: %w", err), makeErr)
		}

		return short, urls.ErrURLAlreadyExists
	} else if err != nil {
		return "", err
	}
	logger.GetSugar().Infof("Added URL: %s -> %s", url.GetOriginal(), url.GetShortURL()) 

	return makeShortURL(url.GetShortURL())
}

func (s *URLService) SetBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {	
	result, err := s.repo.AddBatchURL(batchURLs)
	if err != nil {
		return nil, fmt.Errorf("Unable to add batch URLs: %w", err)
	}

	return result, nil
}

func newID() string {
	return rand.Text()[0:8]
}

func makeShortURL(id string) (string, error) {
	short, err := url.JoinPath("http://"+GetRunAddr(), GetBasePrefix(), id)
	if err != nil {
		return "", err
	}
	return short, nil
}
