package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDBPinger struct {
	pingErr error
}

func (m *mockDBPinger) Ping(_ context.Context) error {
	return m.pingErr
}

func TestGetPing(t *testing.T) {
	initTestLogger(t)

	tests := []struct {
		name     string
		db       DBPinger
		wantCode int
	}{
		{
			name:     "ping success",
			db:       &mockDBPinger{pingErr: nil},
			wantCode: http.StatusOK,
		},
		{
			name:     "ping failure",
			db:       &mockDBPinger{pingErr: errors.New("connection refused")},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "nil db",
			db:       nil,
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewDBHandler(tt.db)
			req, err := http.NewRequest(http.MethodGet, "/ping", nil)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			handler.GetPing(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}
