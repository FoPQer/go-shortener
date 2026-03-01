package service

import "encoding/json"

type InputJSONURL struct {
	URL string `json:"url"`
}

type OutputJSONURL struct {
	Result string `json:"result"`
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