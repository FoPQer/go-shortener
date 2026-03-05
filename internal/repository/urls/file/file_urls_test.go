package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempFile(t *testing.T) string {
	tmpDir := t.TempDir()
	return filepath.Join(tmpDir, "test_urls.json")
}

func createTempFileWithData(t *testing.T, urls []*model.Urls) string {
	filePath := createTempFile(t)
	data, err := json.Marshal(urls)
	require.NoError(t, err)
	err = os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)
	return filePath
}

func TestGetUrls_EmptyFile(t *testing.T) {
	filePath := createTempFile(t)
	repo := NewRepository(filePath)

	result := repo.GetUrls()

	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestGetUrls_ValidData(t *testing.T) {
	testUrls := []*model.Urls{
		{Original: "https://example.com", ShortURL: "GJFTZTEQ"},
		{Original: "https://google.com", ShortURL: "NWEOHOB6"},
	}
	filePath := createTempFileWithData(t, testUrls)
	repo := NewRepository(filePath)

	result := repo.GetUrls()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "https://example.com", result[0].GetOriginal())
	assert.Equal(t, "GJFTZTEQ", result[0].GetShortURL())
	assert.Equal(t, "https://google.com", result[1].GetOriginal())
	assert.Equal(t, "NWEOHOB6", result[1].GetShortURL())
}

func TestGetUrls_InvalidJSON(t *testing.T) {
	filePath := createTempFile(t)
	
	err := os.WriteFile(filePath, []byte("invalid json data"), 0644)
	require.NoError(t, err)

	repo := NewRepository(filePath)
	result := repo.GetUrls()

	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestSetUrls(t *testing.T) {
	tests := []struct {
		name     string
		urls     []*model.Urls
		expected int
	}{
		{
			name: "Set multiple URLs",
			urls: []*model.Urls{
				{Original: "https://new1.com", ShortURL: "new1"},
				{Original: "https://new2.com", ShortURL: "new2"},
			},
			expected: 2,
		},
		{
			name:     "Set empty slice",
			urls:     []*model.Urls{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := createTempFile(t)
			repo := NewRepository(filePath)

			repo.SetUrls(tt.urls)

			result := repo.GetUrls()
			assert.Equal(t, tt.expected, len(result))

			if tt.expected > 0 {
				assert.Equal(t, "https://new1.com", result[0].GetOriginal())
				assert.Equal(t, "new1", result[0].GetShortURL())
				assert.Equal(t, "https://new2.com", result[1].GetOriginal())
				assert.Equal(t, "new2", result[1].GetShortURL())
			}
		})
	}
}

func TestGetURLByShortURL(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) string
		shortURL      string
		expectedURL   string
		expectedError error
	}{
		{
			name: "Found",
			setup: func(t *testing.T) string {
				testUrls := []*model.Urls{
					{Original: "https://example.com", ShortURL: "GJFTZTEQ"},
					{Original: "https://google.com", ShortURL: "NWEOHOB6"},
				}
				return createTempFileWithData(t, testUrls)
			},
			shortURL:      "GJFTZTEQ",
			expectedURL:   "https://example.com",
			expectedError: nil,
		},
		{
			name: "Not found",
			setup: func(t *testing.T) string {
				testUrls := []*model.Urls{
					{Original: "https://example.com", ShortURL: "GJFTZTEQ"},
				}
				return createTempFileWithData(t, testUrls)
			},
			shortURL:      "nonexistent",
			expectedURL:   "",
			expectedError: model.ErrBadValueReceive,
		},
		{
			name: "Empty file",
			setup: func(t *testing.T) string {
				return createTempFile(t)
			},
			shortURL:      "any",
			expectedURL:   "",
			expectedError: model.ErrBadValueReceive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup(t)
			repo := NewRepository(filePath)

			original, err := repo.GetURLByShortURL(tt.shortURL)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedURL, original)
		})
	}
}

func TestAddURL(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		actions  func(t *testing.T, repo *FileUrlsRepository)
		validate func(t *testing.T, repo *FileUrlsRepository)
	}{
		{
			name: "Single URL",
			setup: func(t *testing.T) string {
				return createTempFile(t)
			},
			actions: func(t *testing.T, repo *FileUrlsRepository) {
				u, err := repo.AddURL("https://example.com", "GJFTZTEQ")
				require.NoError(t, err)
				assert.NotNil(t, u)
				assert.Equal(t, "https://example.com", u.GetOriginal())
				assert.Equal(t, "GJFTZTEQ", u.GetShortURL())
			},
			validate: func(t *testing.T, repo *FileUrlsRepository) {
				result := repo.GetUrls()
				assert.Equal(t, 1, len(result))
				assert.Equal(t, "https://example.com", result[0].GetOriginal())
				assert.Equal(t, "GJFTZTEQ", result[0].GetShortURL())
			},
		},
		{
			name: "Multiple URLs",
			setup: func(t *testing.T) string {
				return createTempFile(t)
			},
			actions: func(t *testing.T, repo *FileUrlsRepository) {
				repo.AddURL("https://example1.com", "short1")
				repo.AddURL("https://example2.com", "short2")
				repo.AddURL("https://example3.com", "short3")
			},
			validate: func(t *testing.T, repo *FileUrlsRepository) {
				result := repo.GetUrls()
				assert.Equal(t, 3, len(result))
				
				original, err := repo.GetURLByShortURL("short2")
				require.NoError(t, err)
				assert.Equal(t, "https://example2.com", original)
			},
		},
		{
			name: "Add to existing data",
			setup: func(t *testing.T) string {
				testUrls := []*model.Urls{
					{Original: "https://existing.com", ShortURL: "existing"},
				}
				return createTempFileWithData(t, testUrls)
			},
			actions: func(t *testing.T, repo *FileUrlsRepository) {
				u, err := repo.AddURL("https://new.com", "new")
				require.NoError(t, err)
				assert.NotNil(t, u)
			},
			validate: func(t *testing.T, repo *FileUrlsRepository) {
				result := repo.GetUrls()
				assert.Equal(t, 2, len(result))
				assert.Equal(t, "https://existing.com", result[0].GetOriginal())
				assert.Equal(t, "https://new.com", result[1].GetOriginal())
			},
		},
		{
			name: "Persistence",
			setup: func(t *testing.T) string {
				return createTempFile(t)
			},
			actions: func(t *testing.T, repo *FileUrlsRepository) {
				repo.AddURL("https://example.com", "GJFTZTEQ")
			},
			validate: func(t *testing.T, repo *FileUrlsRepository) {
				repo2 := NewRepository(repo.filePath)
				result := repo2.GetUrls()
				
				assert.Equal(t, 1, len(result))
				assert.Equal(t, "https://example.com", result[0].GetOriginal())
				assert.Equal(t, "GJFTZTEQ", result[0].GetShortURL())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup(t)
			repo := NewRepository(filePath)

			tt.actions(t, repo)
			tt.validate(t, repo)
		})
	}
}

func TestSetUrls_Overwrite(t *testing.T) {
	testUrls := []*model.Urls{
		{Original: "https://old.com", ShortURL: "old"},
	}
	filePath := createTempFileWithData(t, testUrls)
	repo := NewRepository(filePath)

	newUrls := []*model.Urls{
		{Original: "https://new1.com", ShortURL: "new1"},
		{Original: "https://new2.com", ShortURL: "new2"},
	}
	repo.SetUrls(newUrls)

	result := repo.GetUrls()
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "https://new1.com", result[0].GetOriginal())
	assert.Equal(t, "new1", result[0].GetShortURL())

	_, err := repo.GetURLByShortURL("old")
	assert.Error(t, err)
}
