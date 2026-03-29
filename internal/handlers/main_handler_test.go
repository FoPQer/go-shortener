package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/model"
	urlMemory "github.com/FoPQer/go-shortener/internal/repository/urls/memory"
	userMemory "github.com/FoPQer/go-shortener/internal/repository/user/memory"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testLoggerOnce sync.Once

func initTestLogger(t *testing.T) {
	t.Helper()

	testLoggerOnce.Do(func() {
		err := logger.InitLogger()
		require.NoError(t, err)
	})
}

func TestGetUrl(t *testing.T) {
	initTestLogger(t)

	type container struct {
		URLService *service.URLService
	}
	cont := &container{URLService: service.NewURLService(urlMemory.NewRepository())}
	handler := NewHandler(cont.URLService, nil, nil)
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
			urls:  &model.Urls{
				Original: "https://priem.mirea.ru/lk/admin/crud/list/user-resources",
				ShortURL: "67KBAWAO",
			},
			value: "",
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name:  "double get",
			urls:  &model.Urls{
				Original: "https://priem.mirea.ru/lk",
				ShortURL: "RVHUL6VG",
			},
			value: "RVHUL6VG/RVHUL6VG",
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cont.URLService.SetUrls([]*model.Urls{tt.urls})
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
	initTestLogger(t)

	type container struct {
		URLService *service.URLService
	}
	cont := &container{URLService: service.NewURLService(urlMemory.NewRepository())}
	handler := NewHandler(cont.URLService, nil, nil)
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
			urls:  model.NewUrls("", ""),
			value: "https://priem.mirea.ru/lk",
			want: want{
				code:        http.StatusCreated,
				isEmptyBody: false,
			},
		},
		{
			name:  "existing original url",
			urls:  model.NewUrls("https://priem.mirea.ru/lk", "RVHUL6VG"),
			value: "https://priem.mirea.ru/lk",
			want: want{
				code:        http.StatusConflict,
				isEmptyBody: false,
			},
		},
		{
			name:  "empty value",
			urls:  model.NewUrls("", ""),
			value: "",
			want: want{
				code:        http.StatusBadRequest,
				isEmptyBody: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cont.URLService.SetUrls([]*model.Urls{tt.urls})
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewBuffer([]byte(tt.value)))
			w := httptest.NewRecorder()
			handler.PostURL(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			if !tt.want.isEmptyBody {
				assert.NotEmpty(t, resBody)
				if tt.want.code == http.StatusConflict {
					assert.Contains(t, string(resBody), tt.urls.GetShortURL())
				}
			}
		})
	}
}

func TestPostURLByJSON(t *testing.T) {
	initTestLogger(t)

	type container struct {
		URLService *service.URLService
		JSONService *service.JSONService
		UserService *service.UserService
	}
	cont := &container{
		URLService: service.NewURLService(urlMemory.NewRepository()),
		JSONService: service.NewJSONService(),
		UserService: service.NewUserService(userMemory.NewRepository()),
	}
	handler := NewHandler(cont.URLService, cont.JSONService, cont.UserService)
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
			urls:  model.NewUrls("", ""),
			body:  `{"url":"https://priem.mirea.ru/lk"}`,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
				isEmptyBody: false,
			},
		},
		{
			name:  "existing original url",
			urls:  model.NewUrls("https://priem.mirea.ru/lk", "RVHUL6VG"),
			body:  `{"url":"https://priem.mirea.ru/lk"}`,
			want: want{
				code:        http.StatusConflict,
				contentType: "application/json",
				isEmptyBody: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cont.URLService.SetUrls([]*model.Urls{tt.urls})
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
				if tt.want.code == http.StatusConflict {
					assert.Contains(t, string(resBody), tt.urls.GetShortURL())
				}
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
