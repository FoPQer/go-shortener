package repository

import (
	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	urls []*model.Urls
)

func InitUrls() {
	urls = append(urls, model.NewUrls())
}

func GetUrls() []*model.Urls {
	return urls
}

func SetUrls(newUrls []*model.Urls) {
	urls = newUrls
}

func GetURLByShortURL(shortURL string) (string, error) {
	for _, u := range urls {
		if u.GetShortURL() == shortURL {
			return u.GetOriginal(), nil
		}
	}
	return "", model.ErrBadValueReceive
}

func AddURL(original, shortURL string) {
	u := model.NewUrls()
	u.SetOriginal(original)
	u.SetShortURL(shortURL)
	urls = append(urls, u)
}
