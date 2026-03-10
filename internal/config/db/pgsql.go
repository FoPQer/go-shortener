package db

import (
	"context"
	"log"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func GetDBConn() *pgx.Conn {
	return db
}

func SetDBConn(conn *pgx.Conn) {
	db = conn
}

func InitPgsql() *pgx.Conn {
	if service.GetDatabaseDSN() == "" {
		log.Println("Connection to database not found")
		return nil
	}
	conn, err := pgx.Connect(context.Background(), service.GetDatabaseDSN())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil
	}
	SetDBConn(conn)
	log.Println("Connected to database successfully")
	if err := runMigrations(); err != nil {
		logger.GetSugar().Errorf("Unable to run migrations: %v\n", err)
		return nil
	}
	
	return conn
}

func runMigrations() error {
	m, err := migrate.New(
		"file://migrations",
		service.GetDatabaseDSN(),
	)
	if err != nil {
		logger.GetSugar().Errorf("Unable to create migration: %v\n", err)
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.GetSugar().Errorf("Unable to apply migration: %v\n", err)
		return err
	}
	return nil
}