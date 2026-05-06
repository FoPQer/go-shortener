package memory

import (
	"context"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrls(t *testing.T) {
	repo := NewRepository()
	repo.AddURL(context.Background(), "https://example.com", "GJFTZTEQ", "user1")
	repo.AddURL(context.Background(), "https://google.com", "NWEOHOB6", "user2")

	result := repo.GetUrls(context.Background())

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "https://example.com", result[0].GetOriginal())
	assert.Equal(t, "GJFTZTEQ", result[0].GetShortURL())
}

func TestSetUrls(t *testing.T) {
	newUrls := []*model.Urls{
		{Original: "https://new1.com", ShortURL: "new1"},
		{Original: "https://new2.com", ShortURL: "new2"},
	}

	repo := NewRepository()
	repo.SetUrls(context.Background(), newUrls)

	assert.Equal(t, newUrls, repo.GetUrls(context.Background()))
	assert.Equal(t, 2, len(repo.GetUrls(context.Background())))
}

func TestGetURLByShortURL_Found(t *testing.T) {
	repo := NewRepository()
	repo.AddURL(context.Background(), "https://example.com", "GJFTZTEQ", "user1")
	repo.AddURL(context.Background(), "https://google.com", "NWEOHOB6", "user2")

	original, err := repo.GetURLByShortURL(context.Background(), "GJFTZTEQ")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", original)
}

func TestGetURLByShortURL_NotFound(t *testing.T) {
	repo := NewRepository()
	repo.AddURL(context.Background(), "https://example.com", "GJFTZTEQ", "user1")

	original, err := repo.GetURLByShortURL(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.ErrorIs(t, err, urls.ErrURLNotFound)
	assert.Equal(t, "", original)
}

func TestGetURLByShortURL_EmptyURLs(t *testing.T) {
	repo := NewRepository()

	original, err := repo.GetURLByShortURL(context.Background(), "any")

	assert.Error(t, err)
	assert.ErrorIs(t, err, urls.ErrURLNotFound)
	assert.Equal(t, "", original)
}

func TestAddURL(t *testing.T) {
	repo := NewRepository()

	u, err := repo.AddURL(context.Background(), "https://example.com", "GJFTZTEQ", "user1")
	require.NoError(t, err)

	assert.Equal(t, 1, len(repo.GetUrls(context.Background())))
	assert.Equal(t, "https://example.com", u.GetOriginal())
	assert.Equal(t, "GJFTZTEQ", u.GetShortURL())

	u2, err := repo.AddURL(context.Background(), "https://google.com", "NWEOHOB6", "user2")
	require.NoError(t, err)

	assert.Equal(t, 2, len(repo.GetUrls(context.Background())))
	assert.Equal(t, "https://google.com", u2.GetOriginal())
	assert.Equal(t, "NWEOHOB6", u2.GetShortURL())
}

func TestAddURL_MultipleAdd(t *testing.T) {
	repo := NewRepository()

	repo.AddURL(context.Background(), "https://example1.com", "short1", "user1")
	repo.AddURL(context.Background(), "https://example2.com", "short2", "user1")
	repo.AddURL(context.Background(), "https://example3.com", "short3", "user1")

	assert.Equal(t, 3, len(repo.GetUrls(context.Background())))

	result, err := repo.GetURLByShortURL(context.Background(), "short2")
	require.NoError(t, err)
	assert.Equal(t, "https://example2.com", result)
}

func TestAddURL_DuplicateOriginalReturnsExisting(t *testing.T) {
	repo := NewRepository()

	first, err := repo.AddURL(context.Background(), "https://example.com", "old-short", "user1")
	require.NoError(t, err)
	require.NotNil(t, first)

	second, err := repo.AddURL(context.Background(), "https://example.com", "new-short", "user2")
	require.Error(t, err)
	require.ErrorIs(t, err, urls.ErrURLAlreadyExists)
	require.NotNil(t, second)

	assert.Equal(t, "old-short", second.GetShortURL())
	assert.Equal(t, "https://example.com", second.GetOriginal())
	assert.Equal(t, 1, len(repo.GetUrls(context.Background())))
}

func TestAddBatchURL(t *testing.T) {
	repo := NewRepository()

	batch := []*model.Urls{
		{Original: "https://example1.com", ShortURL: "short1", UserID: "user1"},
		{Original: "https://example2.com", ShortURL: "short2", UserID: "user1"},
		{Original: "https://example3.com", ShortURL: "short3", UserID: "user1"},
	}

	result, err := repo.AddBatchURL(context.Background(), batch)
	require.NoError(t, err)
	assert.Equal(t, batch, result)
	assert.Equal(t, 3, len(repo.GetUrls(context.Background())))
}
