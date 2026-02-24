package repository

import (
	"encoding/json"
	"os"

	"github.com/FoPQer/go-shortener/internal/model"
)

var (
	urls []*model.Urls
)

func InitUrls(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		urls = make([]*model.Urls, 0)
		return
	}
	if len(data) == 0 {
		urls = make([]*model.Urls, 0)
	} else {
		json.Unmarshal(data, &urls)
	}
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

func AddURL(original, shortURL string) *model.Urls {
	u := model.NewUrls()
	u.SetOriginal(original)
	u.SetShortURL(shortURL)
	urls = append(urls, u)
	return u
}
