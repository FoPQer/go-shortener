package service

import (
	"crypto/rand"
	"net/url"

	"github.com/FoPQer/go-shortener/internal/repository"
)

func newID() string {
	return rand.Text()[0:8]
}

func GetURL(shortURL string) (string, error) {
	url, err := repository.GetURLByShortURL(shortURL)
	if err != nil {
		return "", err
	}

	return url, nil
}

func SetURL(fullURL string) (string, error) {
	shortURL := newID()
	repository.AddURL(fullURL, shortURL)

	target, err := url.JoinPath("http://"+GetRunAddr(), GetBasePrefix(), shortURL)
	if err != nil {
		return "", err
	}

	return target, nil
}
