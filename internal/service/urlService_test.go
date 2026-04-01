package service

import (
	"context"
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

	deleteFn     func(ctx context.Context, shortUrls []string, userID string) error
	getByShortFn func(ctx context.Context, shortURL string) (string, error)
}

func newMockURLRepository() *mockURLRepository {
	return &mockURLRepository{
		deletedByURL: make(map[string]int),
	}
}

func (m *mockURLRepository) GetUrls(ctx context.Context) []*model.Urls { return nil }

func (m *mockURLRepository) SetUrls(ctx context.Context, newUrls []*model.Urls) {}

func (m *mockURLRepository) GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error) {
	return nil, nil
}

func (m *mockURLRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (*model.Urls, error) {
	return nil, nil
}

func (m *mockURLRepository) GetURLByShortURL(ctx context.Context, shortURL string) (string, error) {
	if m.getByShortFn != nil {
		return m.getByShortFn(ctx, shortURL)
	}
	return "", nil
}

func (m *mockURLRepository) AddURL(ctx context.Context, original, shortURL string, userID string) (*model.Urls, error) {
	return model.NewUrls(original, shortURL), nil
}

func (m *mockURLRepository) AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error) {
	return batchURLs, nil
}

func (m *mockURLRepository) DeleteUrls(ctx context.Context, shortUrls []string, userID string) error {
	m.mu.Lock()
	m.deleteCalls++
	for _, u := range shortUrls {
		m.deletedByURL[u]++
	}
	m.mu.Unlock()

	if m.deleteFn != nil {
		return m.deleteFn(ctx, shortUrls, userID)
	}

	return nil
}

func TestURLService_DeleteUrls_EmptyInput(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	err := svc.DeleteUrls(context.Background(), nil, "user-1")
	require.NoError(t, err)

	repo.mu.Lock()
	defer repo.mu.Unlock()
	assert.Equal(t, 0, repo.deleteCalls)
}

func TestURLService_DeleteUrls_Success(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	shortURLs := []string{"a", "b", "c", "d", "e", "f", "g"}
	err := svc.DeleteUrls(context.Background(), shortURLs, "user-1")
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
	repo.deleteFn = func(ctx context.Context, shortUrls []string, userID string) error {
		if len(shortUrls) > 0 && (shortUrls[0] == "bad-1" || shortUrls[0] == "bad-2") {
			return fmt.Errorf("cannot delete %s", shortUrls[0])
		}
		return nil
	}
	svc := NewURLService(repo)

	err := svc.DeleteUrls(context.Background(), []string{"ok", "bad-1", "good", "bad-2"}, "user-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "errors while deleting URLs")
	assert.Contains(t, err.Error(), "cannot delete bad-1")
	assert.Contains(t, err.Error(), "cannot delete bad-2")
}

func TestURLService_GetURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(ctx context.Context, shortURL string) (string, error) {
			return "https://example.com", nil
		}
		svc := NewURLService(repo)

		got, err := svc.GetURL(context.Background(), "abc")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(ctx context.Context, shortURL string) (string, error) {
			return "", urls.ErrURLNotFound
		}
		svc := NewURLService(repo)

		_, err := svc.GetURL(context.Background(), "abc")
		require.Error(t, err)
		assert.ErrorIs(t, err, urls.ErrURLNotFound)
	})

	t.Run("deleted", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(ctx context.Context, shortURL string) (string, error) {
			return "", urls.ErrURLDeleted
		}
		svc := NewURLService(repo)

		_, err := svc.GetURL(context.Background(), "abc")
		require.Error(t, err)
		assert.ErrorIs(t, err, urls.ErrURLDeleted)
	})

	t.Run("generic repository error", func(t *testing.T) {
		repo := newMockURLRepository()
		repo.getByShortFn = func(ctx context.Context, shortURL string) (string, error) {
			return "", errors.New("db down")
		}
		svc := NewURLService(repo)

		_, err := svc.GetURL(context.Background(), "abc")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to get URL")
	})
}

func BenchmarkDeleteUrls(b *testing.B) {
	repo := newMockURLRepository()
	var calls atomic.Int64
	repo.deleteFn = func(ctx context.Context, shortUrls []string, userID string) error {
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
		err := svc.DeleteUrls(context.Background(), shortURLs, "bench-user")
		if err != nil {
			b.Fatalf("DeleteUrls failed: %v", err)
		}
	}

	if calls.Load() == 0 {
		b.Fatalf("expected benchmark delete function to be called")
	}
}
