package db

import (
	"context"
	"log"

	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

func InitPgsql() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), service.GetDatabaseDSN())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	log.Println("Connected to database successfully")

	return conn, nil
}