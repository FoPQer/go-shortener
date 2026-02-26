package repository

import (
	"os"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitUrls_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "urls*test.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	urls = nil

	InitUrls(tmpFile.Name())

	assert.NotNil(t, urls)
	assert.Equal(t, 0, len(urls))
}

func TestInitUrls_NonExistentFile(t *testing.T) {
	urls = nil
	
	InitUrls("/undefined/file.json")
	
	assert.NotNil(t, urls)
	assert.Equal(t, 0, len(urls))
}

func TestInitUrls_WithValidJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "urls*test.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	jsonData := `[{"original_url":"https://example.com","short_url":"GJFTZTEQ"},{"original_url":"https://google.com","short_url":"NWEOHOB6"}]`
	_, err = tmpFile.WriteString(jsonData)
	require.NoError(t, err)
	tmpFile.Close()

	urls = nil

	InitUrls(tmpFile.Name())

	assert.NotNil(t, urls)
	assert.Greater(t, len(urls), 0)
}

func TestGetUrls(t *testing.T) {
	urls = nil
	urls = append(urls, &model.Urls{Original: "https://example.com", ShortURL: "GJFTZTEQ"})
	urls = append(urls, &model.Urls{Original: "https://google.com", ShortURL: "NWEOHOB6"})

	result := GetUrls()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "https://example.com", result[0].GetOriginal())
	assert.Equal(t, "GJFTZTEQ", result[0].GetShortURL())
}

func TestSetUrls(t *testing.T) {
	newUrls := []*model.Urls{
		{Original: "https://new1.com", ShortURL: "new1"},
		{Original: "https://new2.com", ShortURL: "new2"},
	}

	SetUrls(newUrls)

	assert.Equal(t, newUrls, urls)
	assert.Equal(t, 2, len(GetUrls()))
}

func TestGetURLByShortURL_Found(t *testing.T) {
	urls = nil
	urls = append(urls, &model.Urls{Original: "https://example.com", ShortURL: "GJFTZTEQ"})
	urls = append(urls, &model.Urls{Original: "https://google.com", ShortURL: "NWEOHOB6"})

	original, err := GetURLByShortURL("GJFTZTEQ")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", original)
}

func TestGetURLByShortURL_NotFound(t *testing.T) {
	urls = nil
	urls = append(urls, &model.Urls{Original: "https://example.com", ShortURL: "GJFTZTEQ"})

	original, err := GetURLByShortURL("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, model.ErrBadValueReceive, err)
	assert.Equal(t, "", original)
}

func TestGetURLByShortURL_EmptyURLs(t *testing.T) {
	urls = nil

	original, err := GetURLByShortURL("any")

	assert.Error(t, err)
	assert.Equal(t, model.ErrBadValueReceive, err)
	assert.Equal(t, "", original)
}

func TestAddURL(t *testing.T) {
	urls = nil

	u := AddURL("https://example.com", "GJFTZTEQ")

	assert.Equal(t, 1, len(urls))
	assert.Equal(t, "https://example.com", u.GetOriginal())
	assert.Equal(t, "GJFTZTEQ", u.GetShortURL())

	u2 := AddURL("https://google.com", "NWEOHOB6")

	assert.Equal(t, 2, len(urls))
	assert.Equal(t, "https://google.com", u2.GetOriginal())
	assert.Equal(t, "NWEOHOB6", u2.GetShortURL())
}

func TestAddURL_MultipleAdd(t *testing.T) {
	urls = nil

	AddURL("https://example1.com", "short1")
	AddURL("https://example2.com", "short2")
	AddURL("https://example3.com", "short3")

	assert.Equal(t, 3, len(urls))

	result, err := GetURLByShortURL("short2")
	require.NoError(t, err)
	assert.Equal(t, "https://example2.com", result)
}
