package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FoPQer/go-shortener/internal/model"
	repository "github.com/FoPQer/go-shortener/internal/repository/urls"
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

func (r *FileUrlsRepository) GetURLByOriginalURL(originalURL string) (*model.Urls, error) {
	urls := r.GetUrls()
	for _, u := range urls {
		if u.GetOriginal() == originalURL {
			return u, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", repository.ErrBadValueReceive, originalURL)
}


func (r *FileUrlsRepository) GetURLByShortURL(shortURL string) (string, error) {
	urls := r.GetUrls()
	for _, u := range urls {
		if u.GetShortURL() == shortURL {
			return u.GetOriginal(), nil
		}
	}
	return "", fmt.Errorf("%w: %s", repository.ErrBadValueReceive, shortURL)
}

func (r *FileUrlsRepository) AddURL(original, shortURL string) (*model.Urls, error) {
	urls := r.GetUrls()
	
	u := model.NewUrls(original, shortURL)
	
	urls = append(urls, u)
	
	utils.WriteToFile(r.filePath, urls)
	
	return u, nil
}

func (r *FileUrlsRepository) AddBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {
	urls := r.GetUrls()
	urls = append(urls, batchURLs...)
	utils.WriteToFile(r.filePath, urls)

	return batchURLs, nil
}

