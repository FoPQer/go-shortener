package db

import (
	"context"
	"log"

	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func GetDBConn() *pgx.Conn {
	return db
}

func SetDBConn(conn *pgx.Conn) {
	db = conn
}

func InitPgsql() error {
	conn, err := pgx.Connect(context.Background(), service.GetDatabaseDSN())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	SetDBConn(conn)
	log.Println("Connected to database successfully")

	return nil
}