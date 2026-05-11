package utils

import (
	"encoding/json"
	"os"

	"github.com/FoPQer/go-shortener/internal/model"
)

// WriteToFile serializes URL entities and writes them to the specified file path.
func WriteToFile(filePath string, urls []*model.Urls) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	json.NewEncoder(file).Encode(urls)

	return nil
}
