package service

import (
	"crypto/rand"
	"errors"
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
		return "", err
	}

	return url, nil
}

func (s *URLService) SetURL(fullURL string) (string, error) {
	id := newID()
	url, err := s.repo.AddURL(fullURL, id)
	if errors.Is(err, urls.ErrURLAlreadyExists) {
		short, makeErr := makeShortURL(url.GetShortURL())
		if makeErr != nil {
			return "", errors.Join(err, makeErr)
		}

		return short, urls.ErrURLAlreadyExists
	} else if err != nil {
		return "", err
	}
	logger.GetSugar().Infof("Added URL: %s -> %s", url.GetOriginal(), url.GetShortURL()) 

	return makeShortURL(url.GetShortURL())
}

func (s *URLService) SetBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {
	var result []*model.Urls
	
	for _, u := range batchURLs {
		url, err := s.repo.AddURL(u.GetOriginal(), u.GetShortURL())
		if err != nil {
			return nil, err
		}
		logger.GetSugar().Infof("Added URL: %s -> %s", url.GetOriginal(), url.GetShortURL())
		result = append(result, url)
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
