package repository

import (
	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	Urls *model.Urls
)

func InitUrls() {
	Urls = model.NewUrls()
}
