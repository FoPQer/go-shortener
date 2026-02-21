package service

import "encoding/json"

type InputJsonURL struct {
	Url string `json:"url"`
}

type OutputJsonURL struct {
	Result string `json:"result"`
}

func GetURLFromJson(input []byte) (string, error) {
	var url InputJsonURL
	if err := json.Unmarshal(input, &url); err != nil {
		return "", err
	}

	return url.Url, nil
}

func SetURLToJson(input string) ([]byte, error) {
	url := OutputJsonURL{Result: input}
	
	result, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	return result, nil
}