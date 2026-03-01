package service

import (
	"crypto/rand"
	"log"
	"net/url"

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
	u := s.repo.AddURL(fullURL, newID())
	log.Printf("Added URL: %s -> %s", u.GetOriginal(), u.GetShortURL())
	target, err := url.JoinPath("http://"+GetRunAddr(), GetBasePrefix(), u.GetShortURL())
	if err != nil {
		return "", err
	}

	return target, nil
}

func newID() string {
	return rand.Text()[0:8]
}
