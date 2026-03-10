package db

import (
	"context"
	"log"
	"time"

	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/jackc/pgx/v5"
)

type DBUrlsRepository struct {
	conn *pgx.Conn
	tx pgx.Tx
}

func NewRepository(conn *pgx.Conn) *DBUrlsRepository {
	return &DBUrlsRepository{
		conn: conn,
		tx:   nil,
	}
}

func (r *DBUrlsRepository) BeginTransaction(ctx context.Context) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	r.tx = tx
	return nil
}

func (r *DBUrlsRepository) CommitTransaction(ctx context.Context) error {
	if r.tx == nil {
		return nil
	}
	err := r.tx.Commit(ctx)
	if err != nil {
		return err
	}
	r.tx = nil
	return nil
}

func (r *DBUrlsRepository) RollbackTransaction(ctx context.Context) error {
	if r.tx == nil {
		return nil
	}
	err := r.tx.Rollback(ctx)
	if err != nil {
		return err
	}
	r.tx = nil
	return nil
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
			"INSERT INTO urls (original_url, short_url) VALUES ($1, $2)", 
			u.GetOriginal(), 
			u.GetShortURL(),
		)
		if err != nil {
			continue
		}
	}
	<-ctx.Done()
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

	query := "INSERT INTO urls (original_url, short_url) VALUES ($1, $2)"

	if r.tx != nil {
		_, err := r.tx.Exec(
			ctx,
			query,
			original,
			shortURL,
		)
		if err != nil {
			log.Printf("Error while adding url: %s", err.Error())
			r.tx.Rollback(ctx)
			return nil, err
		}
	} else {
		result, err := r.conn.Exec(
			ctx,
			query,
			original,
			shortURL,
		)
		if err != nil {
			log.Printf("Error while adding url: %s", err.Error())
			return nil, err
		}

		<-ctx.Done()
		log.Printf("Inserted %d row(s)", result.RowsAffected())
	}

	return model.NewUrls(original, shortURL), nil
}

