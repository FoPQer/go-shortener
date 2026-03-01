package file

import (
	"encoding/json"
	"os"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/utils"
)

type FileUrlsRepository struct {
	filePath string
}

func NewRepository(filePath string) *FileUrlsRepository {
	return &FileUrlsRepository{
		filePath: filePath,
	}
}

func (r *FileUrlsRepository) GetUrls() []*model.Urls {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return make([]*model.Urls, 0)
	}

	var urls []*model.Urls
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return make([]*model.Urls, 0)
	}

	return urls
}

func (r *FileUrlsRepository) SetUrls(newUrls []*model.Urls) {
	utils.WriteToFile(r.filePath, newUrls)
}

func (r *FileUrlsRepository) GetURLByShortURL(shortURL string) (string, error) {
	urls := r.GetUrls()
	for _, u := range urls {
		if u.GetShortURL() == shortURL {
			return u.GetOriginal(), nil
		}
	}
	return "", model.ErrBadValueReceive
}

func (r *FileUrlsRepository) AddURL(original, shortURL string) *model.Urls {
	urls := r.GetUrls()
	
	u := model.NewUrls(original, shortURL)
	
	urls = append(urls, u)
	
	utils.WriteToFile(r.filePath, urls)
	
	return u
}
