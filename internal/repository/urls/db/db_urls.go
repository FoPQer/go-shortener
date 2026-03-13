package db

import (
	"context"
	"errors"
	"time"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBUrlsRepository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *DBUrlsRepository {
	return &DBUrlsRepository{
		conn: conn,
	}
}

func (r *DBUrlsRepository) GetUrls() []*model.Urls {
	urls := make([]*model.Urls, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.conn.Query(
		ctx, 
		"SELECT original_url, short_url FROM urls",
	)
	if err != nil {
		return urls
	}
	defer rows.Close()

	for rows.Next() {
		var original, short string
		if err := rows.Scan(&original, &short); err != nil {
			continue
		}
		urls = append(urls, model.NewUrls(original, short))
	}

	<-ctx.Done()
	return urls
}

func (r *DBUrlsRepository) SetUrls(newUrls []*model.Urls) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, u := range newUrls {
		_, err := r.conn.Exec(
			ctx, 
			"INSERT INTO urls (original_url, short_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING", 
			u.GetOriginal(), 
			u.GetShortURL(),
		)
		if err != nil {
			continue
		}
	}
	<-ctx.Done()
}

func (r *DBUrlsRepository) GetURLByOriginalURL(originalURL string) (*model.Urls, error) {
	var short string
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.conn.QueryRow(
		ctx, 
		"SELECT short_url FROM urls WHERE original_url = $1", 
		originalURL,
	).Scan(&short)
	if err != nil {
		return nil, err
	}
	
	<-ctx.Done()
	return model.NewUrls(originalURL, short), nil
}


func (r *DBUrlsRepository) GetURLByShortURL(shortURL string) (string, error) {
	var original string
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.conn.QueryRow(
		ctx, 
		"SELECT original_url FROM urls WHERE short_url = $1", 
		shortURL,
	).Scan(&original)
	if err != nil {
		return "", err
	}
	
	<-ctx.Done()
	return original, nil
}

func (r *DBUrlsRepository) AddURL(original, shortURL string) (*model.Urls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO urls (original_url, short_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING"

	result, err := r.conn.Exec(
		ctx,
		query,
		original,
		shortURL,
	)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		url, err := r.GetURLByOriginalURL(original)
		if err != nil {
			return nil, errors.Join(err, urls.ErrURLAlreadyExists)
		}

		return url, urls.ErrURLAlreadyExists
	}

	<-ctx.Done()
	return model.NewUrls(original, shortURL), nil
}

