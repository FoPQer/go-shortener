package service

import (
	"encoding/json"
	"os"

	"github.com/FoPQer/go-shortener/internal/repository"
)

func WriteToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	json.NewEncoder(file).Encode(repository.GetUrls())

	return nil
}