package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrConnNotFound indicates that database DSN is not configured.
	ErrConnNotFound = errors.New("connection to database not found")
	// ErrUnableToConnect indicates that opening a database connection failed.
	ErrUnableToConnect = errors.New("unable to connect to database")
)

// PgxConf stores PostgreSQL connection configuration.
type PgxConf struct {
	DB *pgxpool.Pool
}

// GetDBConn returns the active PostgreSQL connection pool.
func (p *PgxConf) GetDBConn() *pgxpool.Pool {
	return p.DB
}

// SetDBConn sets the active PostgreSQL connection pool.
func (p *PgxConf) SetDBConn(conn *pgxpool.Pool) {
	p.DB = conn
}

// InitPgsql initializes PostgreSQL connection and applies database migrations.
func InitPgsql() (*PgxConf, error) {
	var pgxConf = &PgxConf{}
	if service.GetDatabaseDSN() == "" {
		return pgxConf, ErrConnNotFound
	}
	conn, err := pgxpool.New(context.Background(), service.GetDatabaseDSN())
	if err != nil {
		return pgxConf, ErrUnableToConnect
	}

	logger.GetSugar().Infoln("Connected to database successfully")
	if err := runMigrations(); err != nil {
		return pgxConf, err
	}

	pgxConf.SetDBConn(conn)
	return pgxConf, nil
}

// runMigrations applies all pending migrations from the migrations directory.
func runMigrations() error {
	m, err := migrate.New(
		"file://migrations",
		service.GetDatabaseDSN(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
