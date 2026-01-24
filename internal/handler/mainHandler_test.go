package handler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
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
		urls  model.Urls
		value string
		want  want
	}{
		{
			name:  "alone get",
			urls:  model.Urls{Urls: map[string]string{"RVHUL6VG": "https://priem.mirea.ru/lk"}},
			value: "RVHUL6VG",
			want: want{
				code:     307,
				location: "https://priem.mirea.ru/lk",
			},
		},
		{
			name:  "without value",
			urls:  model.Urls{Urls: map[string]string{"RVHUL6VG": "https://priem.mirea.ru/lk"}},
			value: "",
			want: want{
				code:     400,
				location: "",
			},
		},
		{
			name:  "double get",
			urls:  model.Urls{Urls: map[string]string{"RVHUL6VG": "https://priem.mirea.ru/lk"}},
			value: "RVHUL6VG/RVHUL6VG",
			want: want{
				code:     400,
				location: "",
			},
		},
		{
			name:  "two urls",
			urls:  model.Urls{Urls: map[string]string{"RVHUL6VG": "https://priem.mirea.ru/lk", "RVHUL6TH": "https://priem.mirea.ru/lk/admin"}},
			value: "RVHUL6VG",
			want: want{
				code:     307,
				location: "https://priem.mirea.ru/lk",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository.Urls = &tt.urls
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/"+tt.value, nil)
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
			repository.Urls = tt.urls
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
			t.Logf("id: %s", resBody)
		})
	}
}
