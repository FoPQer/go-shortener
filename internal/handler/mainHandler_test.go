package handler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrl(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name  string
		urls  *model.Urls
		value string
		want  want
	}{
		{
			name:  "without value",
			urls:  &model.Urls{Urls: map[string]string{"RVHUL6VG": "https://priem.mirea.ru/lk"}},
			value: "",
			want: want{
				code:     400,
				location: "",
			},
		},
		{
			name:  "double get",
			urls:  &model.Urls{Urls: map[string]string{"RVHUL6VG": "https://priem.mirea.ru/lk"}},
			value: "RVHUL6VG/RVHUL6VG",
			want: want{
				code:     400,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.SetUrls(tt.urls)
			target, _ := url.JoinPath("http://localhost:8080", tt.value)
			t.Logf("url: %s", target)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			w := httptest.NewRecorder()
			handler.GetURL(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()

			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}

func TestPostUrl(t *testing.T) {
	type want struct {
		code        int
		isEmptyBody bool
	}
	tests := []struct {
		name  string
		urls  *model.Urls
		value string
		want  want
	}{
		{
			name:  "alone set",
			urls:  model.NewUrls(),
			value: "https://priem.mirea.ru/lk",
			want: want{
				code:        201,
				isEmptyBody: false,
			},
		},
		{
			name:  "empty value",
			urls:  model.NewUrls(),
			value: "",
			want: want{
				code:        400,
				isEmptyBody: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.SetUrls(tt.urls)
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", bytes.NewBuffer([]byte(tt.value)))
			w := httptest.NewRecorder()
			handler.PostURL(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			if !tt.want.isEmptyBody {
				assert.NotEmpty(t, resBody)
			}
		})
	}
}

func TestPostURLByJSON(t *testing.T) {
	type want struct {
		code        int
		contentType string
		isEmptyBody bool
	}
	tests := []struct {
		name  string
		urls  *model.Urls
		body  string
		want  want
	}{
		{
			name:  "valid json",
			urls:  model.NewUrls(),
			body:  `{"url":"https://priem.mirea.ru/lk"}`,
			want: want{
				code:        201,
				contentType: "application/json",
				isEmptyBody: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.SetUrls(tt.urls)
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewBuffer([]byte(tt.body)))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.PostURLByJSON(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			if tt.want.isEmptyBody {
				assert.Empty(t, resBody)
			} else {
				assert.NotEmpty(t, resBody)
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
