package service

import (
	"encoding/json"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
)

type InputJSONURL struct {
	URL string `json:"url"`
}

type OutputJSONURL struct {
	Result string `json:"result"`
}

type InputJSONBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type OutputJSONBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type JSONService struct {
}

func NewJSONService() *JSONService {
	return &JSONService{}
}

func (s *JSONService) GetURLFromJSON(input []byte) (string, error) {
	var url InputJSONURL
	if err := json.Unmarshal(input, &url); err != nil {
		return "", err
	}

	return url.URL, nil
}

func (s *JSONService) SetURLToJSON(input string) ([]byte, error) {
	url := OutputJSONURL{Result: input}
	
	result, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *JSONService) GetBatchURLFromJSON(input []byte) ([]*model.Urls, error) {
	var inputUrls []InputJSONBatchURL
	if err := json.Unmarshal(input, &inputUrls); err != nil {
		return nil, err
	}
	if len(inputUrls) == 0 {
		return nil, urls.ErrBadValueReceive
	}

	return s.GetUrlsFromInputBatchJSON(inputUrls)
}

func (s *JSONService) SetBatchURLToJSON(input []*model.Urls) ([]byte, error) {
	output, err := s.GetOutputBatchJSONFromUrls(input)
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *JSONService) GetUrlsFromInputBatchJSON(input []InputJSONBatchURL) ([]*model.Urls, error) {
	urls := make([]*model.Urls, 0, len(input))
	for _, u := range input {
		urls = append(urls, &model.Urls{
			Original:   u.OriginalURL,
			ShortURL:   u.CorrelationID,
		})
	}

	return urls, nil
}

func (s *JSONService) GetOutputBatchJSONFromUrls(urls []*model.Urls) ([]OutputJSONBatchURL, error) {
	output := make([]OutputJSONBatchURL, 0, len(urls))
	for _, u := range urls {
		short, err := MakeShortURL(u.GetShortURL())
		if err != nil {
			return output, err
		}
		output = append(output, OutputJSONBatchURL{
			CorrelationID: u.GetShortURL(),
			ShortURL:      short,
		})
	}

	return output, nil
}