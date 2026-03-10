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
	id := newID()
	u, err := s.repo.AddURL(fullURL, id)
	if err != nil {
		return "", err
	}
	log.Printf("Added URL: %s -> %s", u.GetOriginal(), u.GetShortURL()) 

	return makeShortURL(u.GetShortURL())
}

func (s *URLService) SetBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {
	var result []*model.Urls
	
	for _, u := range batchURLs {
		url, err := s.repo.AddURL(u.GetOriginal(), u.GetShortURL())
		if err != nil {
			return nil, err
		}
		log.Printf("Added URL: %s -> %s", url.GetOriginal(), url.GetShortURL())
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
