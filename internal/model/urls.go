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

func (u *Urls) SetUrl(id string, url string) error {
	u.Urls[id] = url
	return nil
}

func (u *Urls) GetUrl(id string) (string, error) {
	url, ok := u.Urls[id]
	if !ok {
		return "", errors.New("Value not received")
	}
	return url, nil
}
