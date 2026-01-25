package model

import (
	"errors"
)

var (
	ErrBadValueReceive = errors.New("value not received")
	ErrEmptyURLID      = errors.New("empty id to insert")
	ErrEmptyURLURL     = errors.New("empty url to insert")
	ErrIDAlreadyExists = errors.New("id is already exists")
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
	if id == "" {
		return ErrEmptyUrlId
	}
	if url == "" {
		return ErrEmptyUrlUrl
	}
	_, ok := u.Urls[id]
	if ok {
		return ErrIdAlreadyExists
	}
	u.Urls[id] = url
	return nil
}

func (u *Urls) GetURL(id string) (string, error) {
	url, ok := u.Urls[id]
	if !ok {
		return "", ErrBadValueReceive
	}
	return url, nil
}
