package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/jackc/pgx/v5"
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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	for _, u := range newUrls {
		_, err := r.conn.Exec(
			ctx, 
			"INSERT INTO urls (original_url, short_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING", 
			u.GetOriginal(), 
			u.GetShortURL(),
			u.GetUserID(),
		)
		if err != nil {
			continue
		}
	}
	<-ctx.Done()
}

func (r *DBUrlsRepository) GetUrlsByUserID(userID string) ([]*model.Urls, error) {
	urls := make([]*model.Urls, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	rows, err := r.conn.Query(
		ctx, 
		"SELECT original_url, short_url FROM urls WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var original, short string
		if err := rows.Scan(&original, &short); err != nil {
			continue
		}
		urls = append(urls, model.NewUrls(original, short))
	}

	logger.GetSugar().Infof("urls: %v", urls)

	<-ctx.Done()
	return urls, nil
}
func (r *DBUrlsRepository) DeleteUrlsByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := r.conn.Exec(
		ctx, 
		"DELETE FROM urls WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (r *DBUrlsRepository) GetURLByOriginalURL(originalURL string) (*model.Urls, error) {
	var short string
	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
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
	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
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

func (r *DBUrlsRepository) AddURL(original, shortURL, userID string) (*model.Urls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	query := "INSERT INTO urls (original_url, short_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING"

	result, err := r.conn.Exec(
		ctx,
		query,
		original,
		shortURL,
		userID,
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

func (r *DBUrlsRepository) AddBatchURL(batchURLs []*model.Urls) ([]*model.Urls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	results := make([]*model.Urls, 0, len(batchURLs))
	batch := &pgx.Batch{}

	query := "INSERT INTO urls (original_url, short_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING"
	for _, u := range batchURLs {
		batch.Queue(query, u.GetOriginal(), u.GetShortURL(), u.GetUserID())
	}
	
	batchResults := r.conn.SendBatch(ctx, batch)
	defer batchResults.Close()

	for _, u := range batchURLs {
		result, err := batchResults.Exec()
		if err != nil {
			return nil, fmt.Errorf("unable to Exec() batch at URL %s -> %s: %w", u.GetOriginal(), u.GetShortURL(), err)
		}

		var url *model.Urls

		if result.RowsAffected() == 0 {
			url, err = r.GetURLByOriginalURL(u.GetOriginal())
			if err != nil {
				return nil, errors.Join(fmt.Errorf("unable to get URL by original URL %s: %w", u.GetOriginal(), err), urls.ErrURLAlreadyExists)
			}
		} else {
			url = u
		}

		results = append(results, url)
	}

	<-ctx.Done()
	return results, nil
}
