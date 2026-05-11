package service

import (
	"encoding/json"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
)

// InputJSONURL represents a single URL create request payload.
type InputJSONURL struct {
	URL string `json:"url"`
}

// OutputJSONURL represents a single URL create response payload.
type OutputJSONURL struct {
	Result string `json:"result"`
}

// InputJSONBatchURL represents a single item in a batch create request.
type InputJSONBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// OutputJSONBatchURL represents a single item in a batch create response.
type OutputJSONBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// JSONService provides JSON marshalling and unmarshalling helpers for URL APIs.
type JSONService struct {
}

// NewJSONService constructs a new JSONService instance.
func NewJSONService() *JSONService {
	return &JSONService{}
}

// GetURLFromJSON parses a single URL from JSON input.
func (s *JSONService) GetURLFromJSON(input []byte) (string, error) {
	var url InputJSONURL
	if err := json.Unmarshal(input, &url); err != nil {
		return "", err
	}

	return url.URL, nil
}

// SetURLToJSON serializes a single short URL response to JSON.
func (s *JSONService) SetURLToJSON(input string) ([]byte, error) {
	url := OutputJSONURL{Result: input}

	result, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetBatchURLFromJSON parses a batch create payload into URL entities.
//
// Returns urls.ErrBadValueReceive when the input batch is empty.
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

// SetBatchURLToJSON serializes URL entities into a batch JSON response.
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

// GetUrlsFromInputBatchJSON converts batch DTO input to URL entities.
func (s *JSONService) GetUrlsFromInputBatchJSON(input []InputJSONBatchURL) ([]*model.Urls, error) {
	urls := make([]*model.Urls, 0, len(input))
	for _, u := range input {
		urls = append(urls, &model.Urls{
			Original: u.OriginalURL,
			ShortURL: u.CorrelationID,
		})
	}

	return urls, nil
}

// GetOutputBatchJSONFromUrls converts URL entities to batch response DTOs.
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
