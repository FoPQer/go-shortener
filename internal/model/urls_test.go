package model_test

import (
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUrls_GettersSetters(t *testing.T) {
	u := &model.Urls{}
	original := "https://example.com"
	shortURL := "http://localhost:8080/GJFTZTEQ"

	u.SetOriginal(original)
	u.SetShortURL(shortURL)

	assert.Equal(t, original, u.GetOriginal())
	assert.Equal(t, shortURL, u.GetShortURL())
}
