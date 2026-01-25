package model_test

import (
	"testing"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUrls_SetURL(t *testing.T) {
	type values struct {
		id  string
		url string
	}
	tests := []struct {
		name    string
		u       *model.Urls
		values  values
		wantErr error
	}{
		{
			"example",
			model.NewUrls(),
			values{
				"67KBAWAO",
				"https://priem.mirea.ru/lk/admin/crud/list/user-resources",
			},
			nil,
		},
		{
			"empty id",
			model.NewUrls(),
			values{
				"",
				"https://priem.mirea.ru/lk/admin/crud/list/user-resources",
			},
			model.ErrEmptyUrlId,
		},
		{
			"empty url",
			model.NewUrls(),
			values{
				"67KBAWAO",
				"",
			},
			model.ErrEmptyUrlUrl,
		},
		{
			"exist id",
			&model.Urls{
				Urls: map[string]string{"67KBAWAO": "https://priem.mirea.ru/lk/admin/crud/list/user-resources"},
			},
			values{
				"67KBAWAO",
				"https://priem.mirea.ru/lk/admin/crud/list/user-resources",
			},
			model.ErrIdAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.u.SetURL(tt.values.id, tt.values.url)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				_, ok := tt.u.Urls[tt.values.id]
				assert.True(t, ok)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				if tt.wantErr != model.ErrIdAlreadyExists {
					_, ok := tt.u.Urls[tt.values.id]
					assert.False(t, ok)
				}
			}
		})
	}
}

func TestUrls_GetURL(t *testing.T) {
	tests := []struct {
		name    string
		u       *model.Urls
		id      string
		want    string
		wantErr error
	}{
		{
			"example",
			&model.Urls{
				Urls: map[string]string{"67KBAWAO": "https://priem.mirea.ru/lk/admin/crud/list/user-resources"},
			},
			"67KBAWAO",
			"https://priem.mirea.ru/lk/admin/crud/list/user-resources",
			nil,
		},
		{
			"empty id",
			&model.Urls{
				Urls: map[string]string{"67KBAWAO": "https://priem.mirea.ru/lk/admin/crud/list/user-resources"},
			},
			"",
			"",
			model.ErrBadValueReceive,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := tt.u.GetURL(tt.id)
			require.EqualValues(t, err, tt.wantErr)
			assert.Equal(t, tt.want, url)
		})
	}
}
