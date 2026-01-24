package model

import (
	"errors"
)

type Urls struct {
	Urls map[string]string
}

func NewUrls() *Urls {
	return &Urls{
		Urls: make(map[string]string),
	}
}

func (u *Urls) SetURL(id string, url string) error {
	u.Urls[id] = url
	return nil
}

func (u *Urls) GetURL(id string) (string, error) {
	url, ok := u.Urls[id]
	if !ok {
		return "", errors.New("value not received")
	}
	return url, nil
}
