package repository

import (
	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	urls *model.Urls
)

func InitUrls() {
	urls = model.NewUrls()
}

func GetUrls() *model.Urls {
	return urls
}

func SetUrls(newUrls *model.Urls) {
	urls = newUrls
}
