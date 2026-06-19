package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCounter struct {
	value int
	err   error
}

func (m mockCounter) Count(ctx context.Context) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.value, nil
}

func TestStatService_GetStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := NewStatService(mockCounter{value: 5}, mockCounter{value: 2})

		stats, err := svc.GetStats(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 5, stats.TotalURLs)
		assert.Equal(t, 2, stats.TotalUsers)
	})

	t.Run("not configured", func(t *testing.T) {
		svc := NewStatService(nil, nil)

		_, err := svc.GetStats(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not configured")
	})

	t.Run("urls counter error", func(t *testing.T) {
		svc := NewStatService(mockCounter{err: errors.New("boom")}, mockCounter{value: 2})

		_, err := svc.GetStats(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to count urls")
	})

	t.Run("users counter error", func(t *testing.T) {
		svc := NewStatService(mockCounter{value: 5}, mockCounter{err: errors.New("boom")})

		_, err := svc.GetStats(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to count users")
	})
}
