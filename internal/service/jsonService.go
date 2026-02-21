package service

import "encoding/json"

type InputJSONURL struct {
	URL string `json:"url"`
}

type OutputJSONURL struct {
	Result string `json:"result"`
}

func GetURLFromJSON(input []byte) (string, error) {
	var url InputJSONURL
	if err := json.Unmarshal(input, &url); err != nil {
		return "", err
	}

	return url.URL, nil
}

func SetURLToJSON(input string) ([]byte, error) {
	url := OutputJSONURL{Result: input}
	
	result, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	return result, nil
}