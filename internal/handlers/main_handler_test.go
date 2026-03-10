package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls/memory"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrl(t *testing.T) {
	type container struct {
		URLService *service.URLService
	}
	cont := &container{URLService: service.NewURLService(memory.NewRepository())}
	handler := NewHandler(cont.URLService, nil)
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
	type container struct {
		URLService *service.URLService
	}
	cont := &container{URLService: service.NewURLService(memory.NewRepository())}
	handler := NewHandler(cont.URLService, nil)
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
	type container struct {
		URLService *service.URLService
		JSONService *service.JSONService
	}
	cont := &container{
		URLService: service.NewURLService(memory.NewRepository()),
		JSONService: service.NewJSONService(),
	}
	handler := NewHandler(cont.URLService, cont.JSONService)
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
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestPostBatchURLByJSON(t *testing.T) {
	type container struct {
		URLService  *service.URLService
		JSONService *service.JSONService
	}
	cont := &container{
		URLService:  service.NewURLService(memory.NewRepository()),
		JSONService: service.NewJSONService(),
	}
	handler := NewHandler(cont.URLService, cont.JSONService)

	body := `[
		{"correlation_id":"abc12345","original_url":"https://priem.mirea.ru/lk"},
		{"correlation_id":"xyz67890","original_url":"https://example.com/path"}
	]`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/batch", bytes.NewBuffer([]byte(body)))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.PostBatchURLByJSON(w, request)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, resBody)

	var out []service.OutputJSONBatchURL
	require.NoError(t, json.Unmarshal(resBody, &out))
	require.Len(t, out, 2)

	assert.Equal(t, "abc12345", out[0].CorrelationID)
	assert.Equal(t, "https://priem.mirea.ru/lk", out[0].ShortURL)
	assert.Equal(t, "xyz67890", out[1].CorrelationID)
	assert.Equal(t, "https://example.com/path", out[1].ShortURL)
}


