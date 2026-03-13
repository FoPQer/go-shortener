package factory

import (
	"context"

	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/repository/urls"
	"github.com/FoPQer/go-shortener/internal/repository/urls/db"
	"github.com/FoPQer/go-shortener/internal/repository/urls/file"
	"github.com/FoPQer/go-shortener/internal/repository/urls/memory"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryFactory struct {
	conn *pgxpool.Pool
	filePath string
}

func NewRepositoryFactory(conn *pgxpool.Pool, filePath string) *RepositoryFactory {
	return &RepositoryFactory{conn: conn, filePath: filePath}
}

func (f *RepositoryFactory) CreateUrlsRepository(ctx context.Context) (urls.Repository, error) {
	var repo urls.Repository
	if f.conn != nil {
		repo = db.NewRepository(f.conn)
		logger.GetSugar().Infoln("Working with database repository")
	} else if f.filePath != "" {
		repo = file.NewRepository(f.filePath)
		logger.GetSugar().Infoln("Working with file repository")
	} else {
		repo = memory.NewRepository()
		logger.GetSugar().Infoln("Working with memory repository")
	}

	return repo, nil
}