package db

import (
	"context"
	"errors"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConf struct {
	DB *pgxpool.Pool
}

func (p *PgxConf) GetDBConn() *pgxpool.Pool {
	return p.DB
}

func (p *PgxConf) SetDBConn(conn *pgxpool.Pool) {
	p.DB = conn
}

func InitPgsql() (*PgxConf, error) {
	if service.GetDatabaseDSN() == "" {
		return nil, errors.New("connection to database not found")
	}
	conn, err := pgxpool.New(context.Background(), service.GetDatabaseDSN())
	if err != nil {
		return nil, err
	}
	var pgxConf PgxConf
	pgxConf.SetDBConn(conn)
	logger.GetSugar().Infoln("Connected to database successfully")
	if err := runMigrations(); err != nil {
		return nil, err
	}
	
	return &pgxConf, nil
}

func runMigrations() error {
	m, err := migrate.New(
		"file://migrations",
		service.GetDatabaseDSN(),
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}