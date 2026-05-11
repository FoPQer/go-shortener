package db

import (
	"context"
	"errors"
	"fmt"

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

func (r *DBUrlsRepository) GetUrls(ctx context.Context) []*model.Urls {
	urls := make([]*model.Urls, 0)

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

	if err := rows.Err(); err != nil {
		return urls
	}

	return urls
}

func (r *DBUrlsRepository) SetUrls(ctx context.Context, newUrls []*model.Urls) {
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
}

func (r *DBUrlsRepository) GetUrlsByUserID(ctx context.Context, userID string) ([]*model.Urls, error) {
	urls := make([]*model.Urls, 0)

	rows, err := r.conn.Query(
		ctx,
		"SELECT original_url, short_url FROM urls WHERE user_id = $1 AND is_deleted = FALSE",
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (r *DBUrlsRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (*model.Urls, error) {
	var short string

	err := r.conn.QueryRow(
		ctx,
		"SELECT short_url FROM urls WHERE original_url = $1 AND is_deleted = FALSE",
		originalURL,
	).Scan(&short)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error find by original URL %s: %w", originalURL, urls.ErrURLNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error find by original URL %s: %w", originalURL, urls.ErrBadValueReceive)
	}

	return model.NewUrls(originalURL, short), nil
}

func (r *DBUrlsRepository) GetURLByShortURL(ctx context.Context, shortURL string) (string, error) {
	var original string
	var isDeleted bool

	err := r.conn.QueryRow(
		ctx,
		"SELECT original_url, is_deleted FROM urls WHERE short_url = $1",
		shortURL,
	).Scan(&original, &isDeleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("error find by short URL %s: %w", shortURL, urls.ErrURLNotFound)
	} else if err != nil {
		return "", fmt.Errorf("error find by short URL %s: %w", shortURL, urls.ErrBadValueReceive)
	}

	if isDeleted {
		return "", urls.ErrURLDeleted
	}
	return original, nil
}

func (r *DBUrlsRepository) AddURL(ctx context.Context, original, shortURL string, userID string) (*model.Urls, error) {
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
		url, err := r.GetURLByOriginalURL(ctx, original)
		if err != nil {
			return nil, errors.Join(err, urls.ErrURLAlreadyExists)
		}

		return url, urls.ErrURLAlreadyExists
	}

	return model.NewUrls(original, shortURL), nil
}

func (r *DBUrlsRepository) AddBatchURL(ctx context.Context, batchURLs []*model.Urls) ([]*model.Urls, error) {
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
			url, err = r.GetURLByOriginalURL(ctx, u.GetOriginal())
			if err != nil {
				return nil, errors.Join(fmt.Errorf("unable to get URL by original URL %s: %w", u.GetOriginal(), err), urls.ErrURLAlreadyExists)
			}
		} else {
			url = u
		}

		results = append(results, url)
	}

	return results, nil
}

func (r *DBUrlsRepository) DeleteUrls(ctx context.Context, shortUrls []string, userID string) error {
	_, err := r.conn.Exec(
		ctx,
		"UPDATE urls SET is_deleted = TRUE WHERE short_url = ANY($1) AND user_id = $2",
		shortUrls,
		userID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("error deleting URLs: %v for user %s: %w", shortUrls, userID, urls.ErrURLNotFound)
	} else if err != nil {
		return fmt.Errorf("error while deleting urls: %v: %w", shortUrls, err)
	}

	return nil
}
