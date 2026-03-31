package service

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockURLRepository struct {
	mu           sync.Mutex
	deleteCalls  int
	deletedByURL map[string]int

	deleteFn     func(shortUrls []string, userID string) error
	getByShortFn func(shortURL string) (string, error)
}

func newMockURLRepository() *mockURLRepository {
	return &mockURLRepository{
		deletedByURL: make(map[string]int),
	}
}

func (m *mockURLRepository) GetUrls() []*model.Urls { return nil }

func (m *mockURLRepository) SetUrls(newUrls []*model.Urls) {}

func (m *mockURLRepository) GetUrlsByUserID(userID string) ([]*model.Urls, error) {
	return nil, nil
}

func (m *mockURLRepository) GetURLByOriginalURL(originalURL string) (*model.Urls, error) {
	return nil, nil
}

func (m *mockURLRepository) GetURLByShortURL(shortURL string) (string, error) {
	if m.getByShortFn != nil {
		return m.getByShortFn(shortURL)
	}
	return "", nil
}

func (m *mockURLRepository) AddURL(original, shortURL string, userID string) (*model.Urls, error) {
	return model.NewUrls(original, shortURL), nil
}

func (m *mockURLRepository) AddBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {
	return batchURLs, nil
}

func (m *mockURLRepository) DeleteUrls(shortUrls []string, userID string) error {
	m.mu.Lock()
	m.deleteCalls++
	for _, u := range shortUrls {
		m.deletedByURL[u]++
	}
	m.mu.Unlock()

	if m.deleteFn != nil {
		return m.deleteFn(shortUrls, userID)
	}

	return nil
}

func TestURLService_DeleteUrls_EmptyInput(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	err := svc.DeleteUrls(nil, "user-1")
	require.NoError(t, err)

	repo.mu.Lock()
	defer repo.mu.Unlock()
	assert.Equal(t, 0, repo.deleteCalls)
}

func TestURLService_DeleteUrls_Success(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	shortURLs := []string{"a", "b", "c", "d", "e", "f", "g"}
	err := svc.DeleteUrls(shortURLs, "user-1")
	require.NoError(t, err)

	repo.mu.Lock()
	defer repo.mu.Unlock()
	assert.Equal(t, len(shortURLs), repo.deleteCalls)
	for _, u := range shortURLs {
		assert.Equal(t, 1, repo.deletedByURL[u])
	}
}

func TestURLService_DeleteUrls_ReturnsAggregatedErrors(t *testing.T) {
	repo := newMockURLRepository()
	repo.deleteFn = func(shortUrls []string, userID string) error {
		if len(shortUrls) > 0 && (shortUrls[0] == "bad-1" || shortUrls[0] == "bad-2") {
			return fmt.Errorf("cannot delete %s", shortUrls[0])
		}
		return nil
	}
	svc := NewURLService(repo)

	err := svc.DeleteUrls([]string{"ok", "bad-1", "good", "bad-2"}, "user-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "errors while deleting URLs")
	assert.Contains(t, err.Error(), "cannot delete bad-1")
	assert.Contains(t, err.Error(), "cannot delete bad-2")
}

func TestURLService_GetURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(shortURL string) (string, error) {
			return "https://example.com", nil
		}
		svc := NewURLService(repo)

		got, err := svc.GetURL("abc")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(shortURL string) (string, error) {
			return "", urls.ErrURLNotFound
		}
		svc := NewURLService(repo)

		_, err := svc.GetURL("abc")
		require.Error(t, err)
		assert.ErrorIs(t, err, urls.ErrURLNotFound)
	})

	t.Run("deleted", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(shortURL string) (string, error) {
			return "", urls.ErrURLDeleted
		}
		svc := NewURLService(repo)

		_, err := svc.GetURL("abc")
		require.Error(t, err)
		assert.ErrorIs(t, err, urls.ErrURLDeleted)
	})

	t.Run("generic repository error", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(shortURL string) (string, error) {
			return "", errors.New("db down")
		}
		svc := NewURLService(repo)

		_, err := svc.GetURL("abc")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to get URL")
	})
}

func BenchmarkDeleteUrls(b *testing.B) {
	repo := newMockURLRepository()
	var calls atomic.Int64
	repo.deleteFn = func(shortUrls []string, userID string) error {
		calls.Add(1)
		return nil
	}
	svc := NewURLService(repo)

	shortURLs := make([]string, 1000)
	for i := range shortURLs {
		shortURLs[i] = fmt.Sprintf("short-%d", i)
	}

	b.ReportAllocs()
	
	for b.Loop() {
		err := svc.DeleteUrls(shortURLs, "bench-user")
		if err != nil {
			b.Fatalf("DeleteUrls failed: %v", err)
		}
	}

	if calls.Load() == 0 {
		b.Fatalf("expected benchmark delete function to be called")
	}
}
