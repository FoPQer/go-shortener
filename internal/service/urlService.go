package service

import (
	"crypto/rand"
	"net/url"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/repository"
)

func NewID() string {
	return rand.Text()[0:8]
}

func GetURL(id string) (string, error) {
	urls := repository.GetUrls()

	url, err := urls.GetURL(id)
	if err != nil {
		return "", err
	}

	return url, nil
}

func SetURL(fullURL string) (string, error) {
	urls := repository.GetUrls()

	id := NewID()
	if err := urls.SetURL(id, fullURL); err != nil {
		return "", err
	}

	target, err := url.JoinPath("http://"+flags.GetFlagRunAddr(), flags.GetFlagBasePrefix(), id)
	if err != nil {
		return "", err
	}

	return target, nil
}
